package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"order-service/config"
	"order-service/internal/adapter/repository"
	"order-service/internal/adapter/repository/cache"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/service/transaction"
	"order-service/utils"
	"order-service/utils/conv"
	"strings"
	"sync"

	"github.com/labstack/gommon/log"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, req entity.OrderEntity, userData string) (int64, string, error)
	GetAllOrders(ctx context.Context, query entity.OrderQueryString, userData string) ([]entity.OrderEntity, int64, int64, error)
	GetAllOrdersAdmin(ctx context.Context, query entity.OrderQueryString, userData string) ([]entity.OrderEntity, int64, int64, error)
	GetBatchOrders(ctx context.Context, orderIds []int64, jwtUserData entity.JwtUserData, userData string) ([]entity.OrderEntity, error)
	GetOrderById(ctx context.Context, orderId int64, userId int64, userData string) (*entity.OrderEntity, error)
	GetOrderByIdAdmin(ctx context.Context, orderId int64, userData string) (*entity.OrderEntity, error)
	GetOrderByOrderCode(ctx context.Context, orderCode string, jwtUserData entity.JwtUserData, userData string) (*entity.OrderEntity, error)
	GetOrderIdByOrderCodePublic(ctx context.Context, orderCode string) (int64, error)
	UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, userData string) error

	getAllProductsUsersAdmin(ctx context.Context, orders []entity.OrderEntity, userData string) error
	getProductsUserByIdAdmin(ctx context.Context, order *entity.OrderEntity, userData string) error
	getAllProductsUser(ctx context.Context, orders []entity.OrderEntity, userData string) error
	getProductsUserById(ctx context.Context, order *entity.OrderEntity, userData string) error
}

type orderService struct {
	repo        repository.OrderRepositoryInterface
	repoOutbox  repository.OutboxEventInterface
	repoElastic repository.ElasticRepositoryInterface
	cacheOrder  cache.OrderCacheInterface
	txManager   transaction.TransactionManager
	httpService HttpServiceInterface
	cfg         *config.Config
	logger      *log.Logger
}

func NewOrderService(cfg *config.Config, repo repository.OrderRepositoryInterface, repoOutbox repository.OutboxEventInterface, repoElastic repository.ElasticRepositoryInterface, cacheOrder cache.OrderCacheInterface, httpService HttpServiceInterface, txManager transaction.TransactionManager, logger *log.Logger) OrderServiceInterface {
	return &orderService{
		cfg:         cfg,
		repo:        repo,
		httpService: httpService,
		logger:      logger,
		repoOutbox:  repoOutbox,
		cacheOrder:  cacheOrder,
		txManager:   txManager,
		repoElastic: repoElastic,
	}
}

// GetOrderIdByOrderCodePublic implements [OrderServiceInterface].
func (o *orderService) GetOrderIdByOrderCodePublic(ctx context.Context, orderCode string) (int64, error) {
	var orderId int64

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderEntity, err := o.cacheOrder.GetOrderByOrderCode(txCtx, orderCode, 0)
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

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		switch strings.ToLower(jwtUserData.RoleName) {
		case "customer": // requested by customer
			orderEntities, err := o.repo.GetBatchOrders(txCtx, orderIds, jwtUserData.UserID)
			if err != nil {
				return err
			}

			if err := o.getAllProductsUser(txCtx, orderEntities, userData); err != nil {
				if err.Error() == utils.DATA_NOT_FOUND {
					err := errors.New(utils.RELATION_DATA_NOT_FOUND)
					return err
				}
				return err
			}

			orders = orderEntities

		default: // requested by admin
			orderEntities, err := o.repo.GetBatchOrders(txCtx, orderIds, 0)
			if err != nil {
				return err
			}

			if err := o.getAllProductsUsersAdmin(txCtx, orderEntities, userData); err != nil {
				if err.Error() == utils.DATA_NOT_FOUND {
					err := errors.New(utils.RELATION_DATA_NOT_FOUND)
					return err
				}
				return err
			}

			orders = orderEntities
		}

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

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		switch strings.ToLower(jwtUserData.RoleName) {
		case "customer":
			orderEntity, err := o.cacheOrder.GetOrderByOrderCode(txCtx, orderCode, jwtUserData.UserID)
			if err != nil {
				return err
			}

			if err := o.getProductsUserById(txCtx, orderEntity, userData); err != nil {
				if err.Error() == utils.DATA_NOT_FOUND {
					err := errors.New(utils.RELATION_DATA_NOT_FOUND)
					return err
				}
				return err
			}

			order = orderEntity

		default:
			orderEntity, err := o.cacheOrder.GetOrderByOrderCode(txCtx, orderCode, 0)
			if err != nil {
				return err
			}

			if err := o.getProductsUserByIdAdmin(txCtx, orderEntity, userData); err != nil {
				if err.Error() == utils.DATA_NOT_FOUND {
					err := errors.New(utils.RELATION_DATA_NOT_FOUND)
					return err
				}
				return err
			}

			order = orderEntity
		}

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderByOrderCode: %v", err)
		return nil, err
	}

	return order, nil
}

// UpdateOrderStatus implements [OrderServiceInterface].
func (o *orderService) UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, userData string) error {
	var (
		publishEmailUpdateStatus   = o.cfg.PublisherName.EmailUpdateOrderStatus
		publishElasticUpdateStatus = o.cfg.PublisherName.OrderUpdateStatus
		outboxEventEntities        []entity.OutboxEventEntity
	)

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderEntity, err := o.cacheOrder.GetOrderById(txCtx, req.ID, 0)
		if err != nil {
			return err
		}

		statusReq := strings.ToUpper(req.Status)

		if statusReq != "CANCELED" {
			nextStatus := map[string]string{
				"PENDING":   "CONFIRMED",
				"CONFIRMED": "PROCESS",
				"PROCESS":   "SENDING",
				"SENDING":   "DONE",
			}

			if expected, ok := nextStatus[orderEntity.Status]; ok && statusReq != expected {
				err := errors.New(utils.INVALID_STATUS_TRANSITION)
				return err
			}
		}

		orderEntity.Status = statusReq
		orderEntity.Remarks = req.Remarks

		if err := o.repo.UpdateOrderStatus(txCtx, *orderEntity); err != nil {
			return err
		}

		if err := o.getProductsUserByIdAdmin(txCtx, orderEntity, userData); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err = errors.New(utils.RELATION_DATA_NOT_FOUND)
			}
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

		jsonEmailUpdateStatus, _ := json.Marshal(publishEmailPayload)

		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishEmailUpdateStatus,
			Payload:     string(jsonEmailUpdateStatus),
			AggregateID: fmt.Sprintf("%d", orderEntity.ID),
		})

		publishPushNotifPayload := map[string]any{
			"receiver_email":    "",
			"message":           payloadMessage,
			"subject":           "Update Status Order",
			"type":              "UPDATE_STATUS",
			"receiver_id":       orderEntity.BuyerID,
			"notification_type": "PUSH",
		}

		jsonPushNotif, _ := json.Marshal(publishPushNotifPayload)

		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   utils.NOTIF_PUSH,
			Payload:     string(jsonPushNotif),
			AggregateID: fmt.Sprintf("%d", orderEntity.ID),
		})

		publishElasticPayload := map[string]any{
			"id":      orderEntity.ID,
			"status":  orderEntity.Status,
			"remarks": orderEntity.Remarks,
		}

		jsonElasticUpdate, _ := json.Marshal(publishElasticPayload)

		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishElasticUpdateStatus,
			Payload:     string(jsonElasticUpdate),
			AggregateID: fmt.Sprintf("%d", orderEntity.ID),
		})

		if err := o.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
			return err
		}

		if err := o.cacheOrder.DeleteOrderCache(txCtx, orderEntity.ID, orderEntity.OrderCode); err != nil {
			return err
		}

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] UpdateOrderStatus: %v", err)
		return err
	}

	return nil
}

// CreateOrder implements [OrderServiceInterface].
func (o *orderService) CreateOrder(ctx context.Context, req entity.OrderEntity, userData string) (int64, string, error) {
	var (
		publishUpdateStock  = o.cfg.PublisherName.ProductUpdateStock
		publishOrderCreate  = o.cfg.PublisherName.OrderCreate
		outboxEventEntities []entity.OutboxEventEntity
		orderId             int64
	)

	req.OrderCode = conv.GenerateOrderCode()
	shippingFee := 0
	if strings.ToLower(req.ShippingType) == "delivery" {
		shippingFee = 5000
	}
	req.ShippingFee = int64(shippingFee)
	req.Status = "PENDING"

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderIdCreated, err := o.repo.CreateOrder(txCtx, req)
		if err != nil {
			return err
		}

		orderId = orderIdCreated

		if err := o.cacheOrder.DeleteOrderCache(txCtx, orderIdCreated, ""); err != nil {
			return err
		}

		orderEntity, err := o.cacheOrder.GetOrderById(txCtx, orderIdCreated, 0)
		if err != nil {
			return err
		}

		if err := o.getProductsUserById(txCtx, orderEntity, userData); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err = errors.New(utils.RELATION_DATA_NOT_FOUND)
			}
			return err
		}

		jsonOrderCreate, _ := json.Marshal(orderEntity)

		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishOrderCreate,
			Payload:     string(jsonOrderCreate),
			AggregateID: fmt.Sprintf("%d", orderIdCreated),
		})

		publishPayload := make([]any, 0, len(orderEntity.OrderItems))
		for _, oi := range orderEntity.OrderItems {
			orderItem := map[string]any{
				"product_id": oi.ProductID,
				"quantity":   oi.Quantity,
			}
			publishPayload = append(publishPayload, orderItem)
		}

		jsonUpdateStock, _ := json.Marshal(publishPayload)

		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishUpdateStock,
			Payload:     string(jsonUpdateStock),
			AggregateID: "",
		})

		if err := o.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
			return err
		}

		return nil
	}); err != nil {
		o.cacheOrder.DeleteOrderCache(ctx, orderId, "")
		o.logger.Errorf("[OrderService-1] CreateOrder: %v", err)
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

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderEntities, count, pages, err := o.repo.GetAllOrders(txCtx, query)
		if err != nil {
			return err
		}

		if len(orderEntities) == 0 {
			return nil
		}

		if err := o.getAllProductsUser(txCtx, orderEntities, userData); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		orders, countData, totalPages = orderEntities, count, pages

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetAllOrders: %v", err)
		return nil, 0, 0, err
	}

	return orders, countData, totalPages, nil
}

// GetOrderById implements [OrderServiceInterface].
func (o *orderService) GetOrderById(ctx context.Context, orderId int64, userId int64, userData string) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderEntity, err := o.cacheOrder.GetOrderById(txCtx, orderId, userId)
		if err != nil {
			return err
		}

		if err := o.getProductsUserById(txCtx, orderEntity, userData); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		order = orderEntity

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderById: %v", err)
		return nil, err
	}

	return order, nil
}

// GetOrderByIdAdmin implements [OrderServiceInterface].
func (o *orderService) GetOrderByIdAdmin(ctx context.Context, orderId int64, userData string) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderEntity, err := o.cacheOrder.GetOrderById(txCtx, orderId, 0)
		if err != nil {
			return err
		}

		if err := o.getProductsUserByIdAdmin(txCtx, orderEntity, userData); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		order = orderEntity

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderByIdAdmin: %v", err)
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

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		orderEntities, count, pages, err := o.repo.GetAllOrders(txCtx, query)
		if err != nil {
			return err
		}

		if len(orderEntities) == 0 {
			return nil
		}

		if err := o.getAllProductsUsersAdmin(txCtx, orderEntities, userData); err != nil {
			if err.Error() == utils.DATA_NOT_FOUND {
				err := errors.New(utils.RELATION_DATA_NOT_FOUND)
				return err
			}
			return err
		}

		orders, countData, totalPages = orderEntities, count, pages

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetAllOrdersAdmin: %v", err)
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
