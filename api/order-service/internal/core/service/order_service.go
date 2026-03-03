package service

import (
	"context"
	"errors"
	"fmt"
	"order-service/config"
	"order-service/internal/adapter/repository"
	"order-service/internal/core/domain/entity"
	"order-service/utils"
	"order-service/utils/conv"
	"strings"
	"sync"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type OrderServiceInterface interface {
	GetAllOrdersAdmin(ctx context.Context, query entity.OrderQueryString, userData string) ([]entity.OrderEntity, int64, int64, error)
	GetBatchOrders(ctx context.Context, orderIds []int64, jwtUserData entity.JwtUserData, userData string) ([]entity.OrderEntity, error)
	GetOrderByIdAdmin(ctx context.Context, orderId int64, userData string) (*entity.OrderEntity, error)
	UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, userData string) error
	GetOrderByOrderCode(ctx context.Context, orderCode string, jwtUserData entity.JwtUserData, userData string) (*entity.OrderEntity, error)

	CreateOrder(ctx context.Context, req entity.OrderEntity, userData string) (int64, string, error)
	GetAllOrders(ctx context.Context, query entity.OrderQueryString, userData string) ([]entity.OrderEntity, int64, int64, error)
	GetOrderById(ctx context.Context, orderId int64, userId int64, userData string) (*entity.OrderEntity, error)

	GetOrderIdByOrderCodePublic(ctx context.Context, orderCode string) (int64, error)

	getAllProductsUsersAdmin(ctx context.Context, orders []entity.OrderEntity, userData string) error
	getProductsUserByIdAdmin(ctx context.Context, order *entity.OrderEntity, userData string) error
	getAllProductsUser(ctx context.Context, orders []entity.OrderEntity, userData string) error
	getProductsUserById(ctx context.Context, order *entity.OrderEntity, userData string) error
}

type orderService struct {
	repo        repository.OrderRepositoryInterface
	repoOutbox  repository.OutboxEventInterface
	repoElastic repository.ElasticRepositoryInterface
	httpService HttpServiceInterface
	db          *gorm.DB
	logger      *log.Logger
}

// GetOrderIdByOrderCodePublic implements [OrderServiceInterface].
func (o *orderService) GetOrderIdByOrderCodePublic(ctx context.Context, orderCode string) (int64, error) {
	var orderId int64

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		orderEntity, err := o.repo.GetOrderByOrderCode(ctx, orderCode, 0, tx)
		if err != nil {
			return err
		}

		orderId = orderEntity.ID

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderIdByOrderCodePublic: %v", err)
		return 0, err
	}

	return orderId, nil
}

// GetBatchOrders implements [OrderServiceInterface].
func (o *orderService) GetBatchOrders(ctx context.Context, orderIds []int64, jwtUserData entity.JwtUserData, userData string) ([]entity.OrderEntity, error) {
	var orders []entity.OrderEntity

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		if strings.ToLower(jwtUserData.RoleName) == "customer" {
			orderEntities, err := o.repo.GetBatchOrders(ctx, orderIds, jwtUserData.UserID, tx)
			if err != nil {
				return err
			}

			if err := o.getAllProductsUser(ctx, orderEntities, userData); err != nil {
				return err
			}

			orders = orderEntities

			return nil
		}
		orderEntities, err := o.repo.GetBatchOrders(ctx, orderIds, 0, tx)
		if err != nil {
			return err
		}

		if err := o.getAllProductsUsersAdmin(ctx, orderEntities, userData); err != nil {
			return err
		}

		orders = orderEntities

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetBatchOrders: %v", err)
		return nil, err
	}

	return orders, nil
}

// GetOrderByOrderCode implements [OrderServiceInterface].
func (o *orderService) GetOrderByOrderCode(ctx context.Context, orderCode string, jwtUserData entity.JwtUserData, userData string) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		if strings.ToLower(jwtUserData.RoleName) == "customer" {
			orderEntity, err := o.repo.GetOrderByOrderCode(ctx, orderCode, jwtUserData.UserID, tx)
			if err != nil {
				return err
			}

			if err := o.getProductsUserById(ctx, orderEntity, userData); err != nil {
				return err
			}

			order = orderEntity

			return nil
		}

		orderEntity, err := o.repo.GetOrderByOrderCode(ctx, orderCode, 0, tx)
		if err != nil {
			return err
		}

		if err := o.getProductsUserByIdAdmin(ctx, orderEntity, userData); err != nil {
			return err
		}

		order = orderEntity

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderByOrderCode: %v", err.Error())
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatus implements [OrderServiceInterface].
func (o *orderService) UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, userData string) error {
	publishEmailUpdateStatus := config.NewConfig().PublisherName.EmailUpdateOrderStatus
	publishElasticUpdateStatus := config.NewConfig().PublisherName.OrderUpdateStatus

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		orderEntity, err := o.repo.GetOrderById(ctx, req.ID, 0, tx)
		if err != nil {
			return err
		}

		status := strings.ToLower(orderEntity.Status)
		statusReq := strings.ToLower(req.Status)

		if statusReq != "cancelled" {
			err := errors.New(utils.INVALID_STATUS_TRANSITION)

			nextStatus := map[string]string{
				"pending":   "confirmed",
				"confirmed": "process",
				"process":   "sending",
				"sending":   "done",
			}

			if expected, ok := nextStatus[status]; ok && statusReq != expected {
				return err
			}
		}

		orderEntity.Status = strings.ToUpper(req.Status)
		orderEntity.Remarks = req.Remarks

		if err := o.repo.UpdateOrderStatus(ctx, *orderEntity, tx); err != nil {
			return err
		}

		if err := o.getProductsUserByIdAdmin(ctx, orderEntity, userData); err != nil {
			return err
		}

		payloadMessage := fmt.Sprintf("Hello,\n\nYour order with ID %s has been updated with status: %s.\n\nThank you for shopping with us!", orderEntity.OrderCode, orderEntity.Status)

		publishEmailPayload := map[string]any{
			"receiver_email":    orderEntity.BuyerEmail,
			"message":           payloadMessage,
			"subject":           "Update Status Order",
			"type":              "UPDATE_STATUS",
			"receiver_id":       orderEntity.BuyerID,
			"notification_type": "EMAIL",
		}
		if err := o.repoOutbox.CreateEvent(ctx, publishEmailUpdateStatus, publishEmailPayload, &orderEntity.ID, tx); err != nil {
			return err
		}

		publishPushNotifPayload := map[string]any{
			"receiver_email":    "",
			"message":           payloadMessage,
			"subject":           "Update Status Order",
			"type":              "UPDATE_STATUS",
			"receiver_id":       orderEntity.BuyerID,
			"notification_type": "PUSH",
		}
		if err := o.repoOutbox.CreateEvent(ctx, utils.NOTIF_PUSH, publishPushNotifPayload, &orderEntity.ID, tx); err != nil {
			return err
		}

		publishElasticPayload := map[string]any{
			"id":      orderEntity.ID,
			"status":  orderEntity.Status,
			"remarks": orderEntity.Remarks,
		}
		if err := o.repoOutbox.CreateEvent(ctx, publishElasticUpdateStatus, publishElasticPayload, &orderEntity.ID, tx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] UpdateOrderStatus: %v", err.Error())
		return err
	}

	return nil
}

// CreateOrder implements [OrderServiceInterface].
func (o *orderService) CreateOrder(ctx context.Context, req entity.OrderEntity, userData string) (int64, string, error) {
	publishUpdateStock := config.NewConfig().PublisherName.ProductUpdateStock
	publishOrderCreate := config.NewConfig().PublisherName.OrderCreate

	orderId := int64(0)

	req.OrderCode = conv.GenerateOrderCode()
	shippingFee := 0
	if strings.ToLower(req.ShippingType) == "delivery" {
		shippingFee = 5000
	}
	req.ShippingFee = int64(shippingFee)
	req.Status = "PENDING"

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		id, err := o.repo.CreateOrder(ctx, req, tx)
		if err != nil {
			return err
		}

		orderEntity, err := o.repo.GetOrderById(ctx, id, 0, tx)
		if err != nil {
			return err
		}

		if err := o.getProductsUserById(ctx, orderEntity, userData); err != nil {
			return err
		}

		if err := o.repoOutbox.CreateEvent(ctx, publishOrderCreate, orderEntity, &id, tx); err != nil {
			return err
		}

		publishPayload := make([]any, 0, len(orderEntity.OrderItems))
		for _, oi := range orderEntity.OrderItems {
			orderItem := map[string]any{
				"product_id": oi.ProductID,
				"quantity":   oi.Quantity,
			}
			publishPayload = append(publishPayload, orderItem)
		}

		if err := o.repoOutbox.CreateEvent(ctx, publishUpdateStock, publishPayload, nil, tx); err != nil {
			return err
		}

		orderId = id

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] CreateOrder: %v", err.Error())
		return 0, "", err
	}

	return orderId, req.OrderCode, nil
}

// GetAllOrders implements [OrderServiceInterface].
func (o *orderService) GetAllOrders(ctx context.Context, query entity.OrderQueryString, userData string) ([]entity.OrderEntity, int64, int64, error) {
	orders, countData, totalPages, err := o.repoElastic.SearchOrderElastic(ctx, query)
	if err == nil {
		return orders, countData, totalPages, nil
	}

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		orderEntities, count, pages, err := o.repo.GetAllOrders(ctx, query, tx)
		if err != nil {
			return err
		}

		if len(orderEntities) == 0 {
			return nil
		}

		if err := o.getAllProductsUser(ctx, orderEntities, userData); err != nil {
			return err
		}

		orders, countData, totalPages = orderEntities, count, pages

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetAllOrders: %v", err.Error())
		return nil, 0, 0, err
	}

	return orders, countData, totalPages, nil
}

// GetOrderById implements [OrderServiceInterface].
func (o *orderService) GetOrderById(ctx context.Context, orderId int64, userId int64, userData string) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		orderEntity, err := o.repo.GetOrderById(ctx, orderId, userId, tx)
		if err != nil {
			return err
		}

		if err := o.getProductsUserById(ctx, orderEntity, userData); err != nil {
			return err
		}

		order = orderEntity

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderById: %v", err.Error())
		return nil, err
	}

	return order, nil
}

// GetOrderByIdAdmin implements [OrderServiceInterface].
func (o *orderService) GetOrderByIdAdmin(ctx context.Context, orderId int64, userData string) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		orderEntity, err := o.repo.GetOrderById(ctx, orderId, 0, tx)
		if err != nil {
			return err
		}

		if err := o.getProductsUserByIdAdmin(ctx, orderEntity, userData); err != nil {
			return err
		}

		order = orderEntity

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderByIdAdmin: %v", err.Error())
		return nil, err
	}

	return order, nil
}

// GetAllOrdersAdmin implements [OrderServiceInterface].
func (o *orderService) GetAllOrdersAdmin(ctx context.Context, query entity.OrderQueryString, userData string) ([]entity.OrderEntity, int64, int64, error) {
	orders, countData, totalPages, err := o.repoElastic.SearchOrderElastic(ctx, query)
	if err == nil {
		return orders, countData, totalPages, nil
	}

	if err := o.db.Transaction(func(tx *gorm.DB) error {
		orderEntities, count, pages, err := o.repo.GetAllOrders(ctx, query, tx)
		if err != nil {
			return err
		}

		if len(orderEntities) == 0 {
			return nil
		}

		if err := o.getAllProductsUsersAdmin(ctx, orderEntities, userData); err != nil {
			return err
		}

		orders, countData, totalPages = orderEntities, count, pages

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetAllOrdersAdmin: %v", err.Error())
		return nil, 0, 0, err
	}

	return orders, countData, totalPages, nil
}

// getAllProductsUsersAdmin implements [OrderServiceInterface].
func (o *orderService) getAllProductsUsersAdmin(ctx context.Context, orders []entity.OrderEntity, userData string) error {
	userIds := map[int64]struct{}{}
	productIds := map[int64]struct{}{}
	for _, order := range orders {
		for _, item := range order.OrderItems {
			productIds[item.ProductID] = struct{}{}
		}
		userIds[order.BuyerID] = struct{}{}
	}

	reqProductIds := make([]int64, 0, len(productIds))
	for id := range productIds {
		reqProductIds = append(reqProductIds, id)
	}

	reqUserIds := make([]int64, 0, len(userIds))
	for id := range userIds {
		reqUserIds = append(reqUserIds, id)
	}

	var (
		wg             sync.WaitGroup
		resultProducts map[int64]entity.ProductResponseEntity
		resultUsers    map[int64]entity.UserResponseEntity
		errCh          = make(chan error, 1)
		err            error
	)

	wg.Go(func() {
		resultProducts, err = o.httpService.HttpProductsAllService(reqProductIds, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Go(func() {
		resultUsers, err = o.httpService.HttpUsersAllAdminService(reqUserIds, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	for oIdx, order := range orders {
		if q, ok := resultUsers[order.BuyerID]; ok {
			orders[oIdx].BuyerName = q.Name
		}
		for iIdx, item := range order.OrderItems {
			if p, ok := resultProducts[item.ProductID]; ok {
				orders[oIdx].OrderItems[iIdx].ProductImage = p.ProductImage
			}
		}
	}

	return nil
}

// getProductsUserByIdAdmin implements [OrderServiceInterface].
func (o *orderService) getProductsUserByIdAdmin(ctx context.Context, order *entity.OrderEntity, userData string) error {
	productIds := map[int64]struct{}{}
	for _, item := range order.OrderItems {
		productIds[item.ProductID] = struct{}{}
	}

	reqProductIds := make([]int64, 0, len(productIds))
	for id := range productIds {
		reqProductIds = append(reqProductIds, id)
	}

	reqUserId := order.BuyerID

	var (
		wg             sync.WaitGroup
		resultProducts map[int64]entity.ProductResponseEntity
		resultUsers    *entity.UserResponseEntity
		errCh          = make(chan error, 1)
		err            error
	)

	wg.Go(func() {
		resultProducts, err = o.httpService.HttpProductsAllService(reqProductIds, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Go(func() {
		resultUsers, err = o.httpService.HttpUserByIdAdminService(reqUserId, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	order.BuyerName = resultUsers.Name
	order.BuyerEmail = resultUsers.Email
	order.BuyerPhone = resultUsers.Phone
	order.BuyerAddress = resultUsers.Address

	for iIdx, item := range order.OrderItems {
		if p, ok := resultProducts[item.ProductID]; ok {
			order.OrderItems[iIdx].ProductName = p.ProductName
			order.OrderItems[iIdx].ProductImage = p.ProductImage
			order.OrderItems[iIdx].Price = int64(p.SalePrice)
		}
	}

	return nil
}

// getAllProductsUser implements [OrderServiceInterface].
func (o *orderService) getAllProductsUser(ctx context.Context, orders []entity.OrderEntity, userData string) error {
	productIds := map[int64]struct{}{}
	for _, order := range orders {
		for _, item := range order.OrderItems {
			productIds[item.ProductID] = struct{}{}
		}
	}

	reqProductIds := make([]int64, 0, len(productIds))
	for id := range productIds {
		reqProductIds = append(reqProductIds, id)
	}

	var (
		wg             sync.WaitGroup
		resultProducts map[int64]entity.ProductResponseEntity
		resultUsers    *entity.UserResponseEntity
		errCh          = make(chan error, 1)
		err            error
	)

	wg.Go(func() {
		resultProducts, err = o.httpService.HttpProductsAllService(reqProductIds, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Go(func() {
		resultUsers, err = o.httpService.HttpUserByIdService(userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	for oIdx, order := range orders {
		orders[oIdx].BuyerName = resultUsers.Name
		orders[oIdx].BuyerEmail = resultUsers.Email
		orders[oIdx].BuyerPhone = resultUsers.Phone
		orders[oIdx].BuyerAddress = resultUsers.Address
		for iIdx, item := range order.OrderItems {
			if p, ok := resultProducts[item.ProductID]; ok {
				orders[oIdx].OrderItems[iIdx].ProductName = p.ProductName
				orders[oIdx].OrderItems[iIdx].ProductImage = p.ProductImage
				orders[oIdx].OrderItems[iIdx].Price = int64(p.SalePrice)
			}
		}
	}

	return nil
}

// getProductsUserById implements [OrderServiceInterface].
func (o *orderService) getProductsUserById(ctx context.Context, order *entity.OrderEntity, userData string) error {
	productIds := map[int64]struct{}{}
	for _, item := range order.OrderItems {
		productIds[item.ProductID] = struct{}{}
	}

	reqProductIds := make([]int64, 0, len(productIds))
	for id := range productIds {
		reqProductIds = append(reqProductIds, id)
	}

	var (
		wg             sync.WaitGroup
		resultProducts map[int64]entity.ProductResponseEntity
		resultUsers    *entity.UserResponseEntity
		errCh          = make(chan error, 1)
		err            error
	)

	wg.Go(func() {
		resultProducts, err = o.httpService.HttpProductsAllService(reqProductIds, userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Go(func() {
		resultUsers, err = o.httpService.HttpUserByIdService(userData)
		if err != nil {
			errCh <- err
		}
	})

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	order.BuyerName = resultUsers.Name
	order.BuyerEmail = resultUsers.Email
	order.BuyerPhone = resultUsers.Phone
	order.BuyerAddress = resultUsers.Address

	for iIdx, item := range order.OrderItems {
		if p, ok := resultProducts[item.ProductID]; ok {
			order.OrderItems[iIdx].ProductName = p.ProductName
			order.OrderItems[iIdx].ProductImage = p.ProductImage
			order.OrderItems[iIdx].Price = int64(p.SalePrice)
		}
	}

	return nil
}

func NewOrderService(repo repository.OrderRepositoryInterface, repoOutbox repository.OutboxEventInterface, repoElastic repository.ElasticRepositoryInterface, httpService HttpServiceInterface, db *gorm.DB, logger *log.Logger) OrderServiceInterface {
	return &orderService{
		repo:        repo,
		httpService: httpService,
		db:          db,
		logger:      logger,
		repoOutbox:  repoOutbox,
		repoElastic: repoElastic,
	}
}
