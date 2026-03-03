package repository

import (
	"context"
	"errors"
	"math"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/domain/model"
	"order-service/utils"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	GetAllOrders(ctx context.Context, query entity.OrderQueryString, db *gorm.DB) ([]entity.OrderEntity, int64, int64, error)
	GetBatchOrders(ctx context.Context, orderIds []int64, userId int64, db *gorm.DB) ([]entity.OrderEntity, error)
	GetOrderById(ctx context.Context, orderId int64, userId int64, db *gorm.DB) (*entity.OrderEntity, error)
	CreateOrder(ctx context.Context, req entity.OrderEntity, db *gorm.DB) (int64, error)
	UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, db *gorm.DB) error
	DeleteOrder(ctx context.Context, orderId int64, db *gorm.DB) error
	GetOrderByOrderCode(ctx context.Context, orderCode string, userId int64, db *gorm.DB) (*entity.OrderEntity, error)
}

type orderRepository struct {
	logger *log.Logger
}

// GetBatchOrders implements [OrderRepositoryInterface].
func (o *orderRepository) GetBatchOrders(ctx context.Context, orderIds []int64, userId int64, db *gorm.DB) ([]entity.OrderEntity, error) {
	chunkSize := 150

	modelOrders := []model.Order{}

	for i := 0; i < len(orderIds); i += chunkSize {
		end := min(i+chunkSize, len(orderIds))
		batchOrders := []model.Order{}

		sqlMain := db.WithContext(ctx).
			Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
				return db.Select("id", "order_id", "product_id", "quantity")
			}).
			Select("id", "buyer_id", "order_code", "shipping_type").
			Where("id IN ?", orderIds[i:end])

		if userId != 0 {
			sqlMain = sqlMain.Where("buyer_id = ?", userId)
		}

		if err := sqlMain.
			Find(&batchOrders).Error; err != nil {
			o.logger.Errorf("[OrderRepository-1] GetBatchOrders: %v", err)
			return nil, err
		}

		modelOrders = append(modelOrders, batchOrders...)
	}

	if len(modelOrders) == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		o.logger.Errorf("[OrderRepository-2] GetBatchOrders: %v", err)
		return nil, err
	}

	entities := []entity.OrderEntity{}
	for _, val := range modelOrders {
		orderItemEntities := []entity.OrderItemEntity{}
		for _, prd := range val.OrderItems {
			orderItemEntities = append(orderItemEntities, entity.OrderItemEntity{
				ID:        prd.ID,
				ProductID: prd.ProductID,
				Quantity:  prd.Quantity,
			})
		}

		entities = append(entities, entity.OrderEntity{
			ID:           val.ID,
			OrderCode:    val.OrderCode,
			ShippingType: val.ShippingType,
			BuyerID:      val.BuyerID,
			OrderItems:   orderItemEntities,
		})
	}

	return entities, nil
}

// GetOrderByOrderCode implements [OrderRepositoryInterface].
func (o *orderRepository) GetOrderByOrderCode(ctx context.Context, orderCode string, userId int64, db *gorm.DB) (*entity.OrderEntity, error) {
	modelOrder := model.Order{}

	sqlMain := db.WithContext(ctx).
		Where("order_code = ?", orderCode).
		Omit("created_at", "updated_at", "deleted_at").
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_id", "quantity")
		})

	if userId != 0 {
		sqlMain = sqlMain.Where("buyer_id = ?", userId)
	}

	if err := sqlMain.
		First(&modelOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		o.logger.Errorf("[OrderRepository-1] GetOrderByOrderCode: %v", err.Error())
		return nil, err
	}

	orderItemEntities := []entity.OrderItemEntity{}
	for _, item := range modelOrder.OrderItems {
		orderItemEntities = append(orderItemEntities, entity.OrderItemEntity{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &entity.OrderEntity{
		ID:           modelOrder.ID,
		OrderCode:    modelOrder.OrderCode,
		OrderDate:    modelOrder.OrderDate.Format("2006-01-02"),
		OrderTime:    modelOrder.OrderTime,
		Status:       modelOrder.Status,
		BuyerID:      modelOrder.BuyerID,
		TotalAmount:  int64(modelOrder.TotalAmount),
		Remarks:      modelOrder.Remarks,
		ShippingType: modelOrder.ShippingType,
		ShippingFee:  int64(modelOrder.ShippingFee),
		OrderItems:   orderItemEntities,
	}, nil
}

// CreateOrder implements [OrderRepositoryInterface].
func (o *orderRepository) CreateOrder(ctx context.Context, req entity.OrderEntity, db *gorm.DB) (int64, error) {
	orderDate, err := time.Parse("2006-01-02", req.OrderDate)
	if err != nil {
		o.logger.Errorf("[OrderRepository-1] CreateOrder: %v", err.Error())
		return 0, err
	}

	orderItems := []model.OrderItem{}
	for _, item := range req.OrderItems {
		orderItem := model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		orderItems = append(orderItems, orderItem)
	}

	modelOrder := model.Order{
		OrderCode:    req.OrderCode,
		BuyerID:      req.BuyerID,
		OrderDate:    orderDate,
		OrderTime:    req.OrderTime,
		Status:       req.Status,
		TotalAmount:  float64(req.TotalAmount),
		ShippingType: req.ShippingType,
		ShippingFee:  float64(req.ShippingFee),
		Remarks:      req.Remarks,
		OrderItems:   orderItems,
	}

	if err := db.WithContext(ctx).Create(&modelOrder).Error; err != nil {
		o.logger.Errorf("[OrderRepository-2] CreateOrder: %v", err.Error())
		return 0, err
	}

	return modelOrder.ID, nil
}

// DeleteOrder implements [OrderRepositoryInterface].
func (o *orderRepository) DeleteOrder(ctx context.Context, orderId int64, db *gorm.DB) error {
	panic("unimplemented")
}

// GetAllOrders implements [OrderRepositoryInterface].
func (o *orderRepository) GetAllOrders(ctx context.Context, query entity.OrderQueryString, db *gorm.DB) ([]entity.OrderEntity, int64, int64, error) {
	modelOrders := []model.Order{}

	var countData int64

	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_id", "quantity")
		}).
		Where("order_code ILIKE ? OR status ILIKE ?", "%"+query.Search+"%", "%"+query.Status+"%")

	if query.BuyerID != 0 {
		sqlMain = sqlMain.Where("buyer_id = ?", query.BuyerID)
	}

	if err := sqlMain.Model(&modelOrders).Count(&countData).Error; err != nil {
		o.logger.Errorf("[OrderRepository-1] GetAllOrders: %v", err.Error())
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order("id DESC").
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelOrders).Error; err != nil {
		o.logger.Errorf("[OrderRepository-2] GetAllOrders: %v", err.Error())
		return nil, 0, 0, err
	}

	entities := []entity.OrderEntity{}
	for _, val := range modelOrders {
		orderItemEntities := []entity.OrderItemEntity{}
		for _, prd := range val.OrderItems {
			orderItemEntities = append(orderItemEntities, entity.OrderItemEntity{
				ID:        prd.ID,
				ProductID: prd.ProductID,
				Quantity:  prd.Quantity,
			})
		}

		entities = append(entities, entity.OrderEntity{
			ID:          val.ID,
			OrderCode:   val.OrderCode,
			Status:      val.Status,
			OrderDate:   val.OrderDate.Format("2006-01-02"),
			OrderTime:   val.OrderTime,
			TotalAmount: int64(val.TotalAmount),
			BuyerID:     val.BuyerID,
			OrderItems:  orderItemEntities,
		})
	}

	return entities, countData, int64(totalPages), nil
}

// GetOrderById implements [OrderRepositoryInterface].
func (o *orderRepository) GetOrderById(ctx context.Context, orderId int64, userId int64, db *gorm.DB) (*entity.OrderEntity, error) {
	modelOrder := model.Order{}

	sqlMain := db.WithContext(ctx).
		Where("id = ?", orderId).
		Omit("created_at", "updated_at", "deleted_at").
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_id", "quantity")
		})

	if userId != 0 {
		sqlMain = sqlMain.Where("buyer_id = ?", userId)
	}

	if err := sqlMain.First(&modelOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		o.logger.Errorf("[OrderRepository-1] GetOrderById: %v", err.Error())
		return nil, err
	}

	orderItemEntities := []entity.OrderItemEntity{}
	for _, item := range modelOrder.OrderItems {
		orderItemEntities = append(orderItemEntities, entity.OrderItemEntity{
			ID:        item.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &entity.OrderEntity{
		ID:           modelOrder.ID,
		OrderCode:    modelOrder.OrderCode,
		OrderDate:    modelOrder.OrderDate.Format("2006-01-02"),
		OrderTime:    modelOrder.OrderTime,
		Status:       modelOrder.Status,
		BuyerID:      modelOrder.BuyerID,
		TotalAmount:  int64(modelOrder.TotalAmount),
		Remarks:      modelOrder.Remarks,
		ShippingType: modelOrder.ShippingType,
		ShippingFee:  int64(modelOrder.ShippingFee),
		OrderItems:   orderItemEntities,
	}, nil
}

// UpdateOrderStatus implements [OrderRepositoryInterface].
func (o *orderRepository) UpdateOrderStatus(ctx context.Context, req entity.OrderEntity, db *gorm.DB) error {
	modelOrder := model.Order{
		ID:      req.ID,
		Status:  req.Status,
		Remarks: req.Remarks,
	}

	if err := db.WithContext(ctx).Updates(&modelOrder).Error; err != nil {
		o.logger.Errorf("[OrderRepository-1] UpdateOrderStatus: %v", err.Error())
		return err
	}

	return nil
}

func NewOrderRepository(logger *log.Logger) OrderRepositoryInterface {
	return &orderRepository{logger: logger}
}
