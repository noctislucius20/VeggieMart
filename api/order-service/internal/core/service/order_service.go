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
	GetOrderById(ctx context.Context, orderId int64, jwtUserData entity.JwtUserData) (*entity.OrderEntity, error)
	GetOrderByOrderCode(ctx context.Context, orderCode string, jwtUserData entity.JwtUserData) (*entity.OrderEntity, error)
	UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, userData string) error

	getUserHttp(ctx context.Context, order *entity.OrderEntity, userData string) error
	getProductsHttp(ctx context.Context, order *entity.OrderEntity, userData string) error
	updateStockProduct(ctx context.Context, order *entity.OrderEntity, userData string) error
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

// GetOrderByOrderCode implements [OrderServiceInterface].
func (o *orderService) GetOrderByOrderCode(ctx context.Context, orderCode string, jwtUserData entity.JwtUserData) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := o.cacheOrder.GetRoleById(txCtx, jwtUserData.RoleID)
		if err != nil {
			return err
		}

		switch strings.ToLower(roleEntity.Name) {
		case "customer":
			orderEntity, err := o.cacheOrder.GetOrderByOrderCode(txCtx, orderCode, jwtUserData.UserID)
			if err != nil {
				return err
			}

			order = orderEntity

		default:
			orderEntity, err := o.cacheOrder.GetOrderByOrderCode(txCtx, orderCode, 0)
			if err != nil {
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
				return errors.New(utils.INVALID_STATUS_TRANSITION)
			}
		}

		orderEntity.Status = statusReq
		orderEntity.Remarks = req.Remarks

		if err := o.repo.UpdateOrderStatus(txCtx, *orderEntity); err != nil {
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
		publishOrderCreate        = o.cfg.PublisherName.OrderCreate
		publishOrderPaymentCreate = o.cfg.PublisherName.OrderPaymentCreate
		outboxEventEntities       []entity.OutboxEventEntity
		orderId                   int64
		wg                        sync.WaitGroup
		errCh                     = make(chan error, 1)
	)

	req.OrderCode = conv.GenerateOrderCode()
	if strings.ToLower(req.ShippingType) == "delivery" {
		req.ShippingFee = 5000
	}
	req.Status = "PENDING"

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		wg.Go(func() {
			if err := o.getUserHttp(txCtx, &req, userData); err != nil {
				errCh <- err
			}
		})

		wg.Go(func() {
			if err := o.getProductsHttp(txCtx, &req, userData); err != nil {
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

		for _, v := range req.OrderItems {
			req.TotalAmount += v.Price * v.Quantity
		}
		req.TotalAmount += req.ShippingFee

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

		orderEntity.PaymentMethod = req.PaymentMethod

		jsonOrderCreate, _ := json.Marshal(orderEntity)

		// consumed by elastic
		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishOrderCreate,
			Payload:     string(jsonOrderCreate),
			AggregateID: fmt.Sprintf("%d", orderIdCreated),
		})

		// consumed by payment service db
		outboxEventEntities = append(outboxEventEntities, entity.OutboxEventEntity{
			EventType:   publishOrderPaymentCreate,
			Payload:     string(jsonOrderCreate),
			AggregateID: fmt.Sprintf("%d", orderIdCreated),
		})

		if err := o.repoOutbox.CreateBatchEvents(txCtx, outboxEventEntities); err != nil {
			return err
		}

		if err := o.updateStockProduct(txCtx, &req, userData); err != nil {
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

		orders, countData, totalPages = orderEntities, count, pages

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetAllOrders: %v", err)
		return nil, 0, 0, err
	}

	return orders, countData, totalPages, nil
}

// GetOrderById implements [OrderServiceInterface].
func (o *orderService) GetOrderById(ctx context.Context, orderId int64, jwtUserData entity.JwtUserData) (*entity.OrderEntity, error) {
	order := &entity.OrderEntity{}

	if err := o.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := o.cacheOrder.GetRoleById(txCtx, jwtUserData.RoleID)
		if err != nil {
			return err
		}

		switch strings.ToLower(roleEntity.Name) {
		case "customer": // requested by customer
			orderEntity, err := o.cacheOrder.GetOrderById(txCtx, orderId, jwtUserData.UserID)
			if err != nil {
				return err
			}

			order = orderEntity
		default: // requested by admin
			orderEntity, err := o.cacheOrder.GetOrderById(txCtx, orderId, 0)
			if err != nil {
				return err
			}

			order = orderEntity
		}

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetOrderById: %v", err)
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

		orders, countData, totalPages = orderEntities, count, pages

		return nil
	}); err != nil {
		o.logger.Errorf("[OrderService-1] GetAllOrdersAdmin: %v", err)
		return nil, 0, 0, err
	}

	return orders, countData, totalPages, nil
}

// updateStockProduct implements [OrderServiceInterface].
func (o *orderService) updateStockProduct(ctx context.Context, order *entity.OrderEntity, userData string) error {
	if err := o.httpService.HttpUpdateStockProductsService(order.OrderItems, userData); err != nil {
		return err
	}

	return nil
}

// getProductsHttp implements [OrderServiceInterface].
func (o *orderService) getProductsHttp(ctx context.Context, order *entity.OrderEntity, userData string) error {
	productIds := map[int64]struct{}{}
	for _, item := range order.OrderItems {
		productIds[item.ProductID] = struct{}{}
	}

	reqProductIds := make([]int64, 0, len(productIds))
	for id := range productIds {
		reqProductIds = append(reqProductIds, id)
	}

	resultProducts, err := o.httpService.HttpGetProductsService(reqProductIds, userData)
	if err != nil {
		return err
	}

	for iIdx, item := range order.OrderItems {
		if p, ok := resultProducts[item.ProductID]; ok {
			order.OrderItems[iIdx].ProductName = p.ProductName
			order.OrderItems[iIdx].ProductImage = p.ProductImage
			order.OrderItems[iIdx].Price = int64(p.SalePrice)
			order.ProductSnapshots = append(order.ProductSnapshots, entity.ProductSnapshotEntity{
				ProductID:    p.ID,
				Name:         p.ProductName,
				Image:        p.ProductImage,
				RegularPrice: p.RegularPrice,
				SalePrice:    p.SalePrice,
				Weight:       int64(p.Weight),
				Unit:         p.Unit,
			})
		}
	}

	return nil
}

// getUserHttp implements [OrderServiceInterface].
func (o *orderService) getUserHttp(ctx context.Context, order *entity.OrderEntity, userData string) error {
	resultUsers, err := o.httpService.HttpUserByIdService(userData)
	if err != nil {
		return err
	}

	order.BuyerName = resultUsers.Name
	order.BuyerEmail = resultUsers.Email
	order.BuyerPhone = resultUsers.Phone
	order.BuyerAddress = resultUsers.Address
	order.BuyerLat = resultUsers.Lat
	order.BuyerLng = resultUsers.Lng

	return nil
}
