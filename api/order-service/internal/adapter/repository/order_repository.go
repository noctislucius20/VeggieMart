package repository

import (
	"context"
	"errors"
	"math"
	"order-service/internal/core/domain/entity"
	"order-service/internal/core/domain/model"
	"order-service/utils"
	"order-service/utils/conv"
	"strings"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type OrderRepositoryInterface interface {
	GetAllOrders(ctx context.Context, query entity.OrderQueryString) ([]entity.OrderEntity, int64, int64, error)
	GetOrderById(ctx context.Context, orderId int64, userId int64) (*entity.OrderEntity, error)
	CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error)
	UpdateOrderStatus(ctx context.Context, req entity.OrderEntity) error
	DeleteOrder(ctx context.Context, orderId int64) error
	GetOrderByOrderCode(ctx context.Context, orderCode string, userId int64) (*entity.OrderEntity, error)

	getDB(ctx context.Context) *gorm.DB
}

type orderRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewOrderRepository(db *gorm.DB, logger *log.Logger) OrderRepositoryInterface {
	return &orderRepository{db: db, logger: logger}
}

// getDB implements [OrderRepositoryInterface].
func (o *orderRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return o.db
}

// GetOrderByOrderCode implements [OrderRepositoryInterface].
func (o *orderRepository) GetOrderByOrderCode(ctx context.Context, orderCode string, userId int64) (*entity.OrderEntity, error) {
	var (
		db          = o.getDB(ctx)
		modelOrder  model.Order
		orderEntity *entity.OrderEntity
		snapshotMap = make(map[int64]entity.ProductSnapshotEntity)
	)

	sqlMain := db.WithContext(ctx).
		Where("order_code = ?", orderCode).
		Omit("created_at", "updated_at", "deleted_at").
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_id", "quantity")
		}).
		Preload("ProductSnapshots", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "last_used", "regular_price")
		})

	if err := sqlMain.
		First(&modelOrder).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = utils.ErrDataNotFound
		}
		o.logger.Errorf("[OrderRepository] GetOrderByOrderCode: %v", err)
		return nil, err
	}

	if userId != 0 && modelOrder.BuyerID != userId {
		err := utils.ErrAccessForbidden
		o.logger.Errorf("[OrderRepository] GetOrderByOrderCode: %v", err)
		return nil, err
	}

	orderEntity = &entity.OrderEntity{
		ID:           modelOrder.ID,
		OrderCode:    modelOrder.OrderCode,
		OrderDate:    modelOrder.OrderDate.Format("2006-01-02"),
		OrderTime:    modelOrder.OrderTime.Format("15:04:05"),
		Status:       modelOrder.Status,
		BuyerID:      modelOrder.BuyerID,
		TotalAmount:  int64(modelOrder.TotalAmount),
		Remarks:      modelOrder.Remarks,
		ShippingType: modelOrder.ShippingType,
		ShippingFee:  int64(modelOrder.ShippingFee),
		BuyerName:    modelOrder.BuyerName,
		BuyerEmail:   modelOrder.BuyerEmail,
		BuyerPhone:   modelOrder.BuyerPhone,
		BuyerAddress: modelOrder.BuyerAddress,
		BuyerLat:     modelOrder.BuyerLat,
		BuyerLng:     modelOrder.BuyerLng,
	}

	for _, snapshot := range modelOrder.ProductSnapshots {
		snapshotMap[snapshot.ProductID] = entity.ProductSnapshotEntity{
			ProductID: snapshot.ProductID,
			Name:      snapshot.Name,
			Image:     snapshot.Image,
			SalePrice: snapshot.SalePrice,
			Weight:    snapshot.Weight,
			Unit:      snapshot.Unit,
		}
	}

	for _, item := range modelOrder.OrderItems {
		if snapshot, ok := snapshotMap[item.ProductID]; ok {
			orderEntity.OrderItems = append(orderEntity.OrderItems, entity.OrderItemEntity{
				ID:        item.ID,
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,

				ProductName:  snapshot.Name,
				ProductImage: snapshot.Image,
				Price:        int64(snapshot.SalePrice),
			})
		}
	}

	return orderEntity, nil
}

// CreateOrder implements [OrderRepositoryInterface].
func (o *orderRepository) CreateOrder(ctx context.Context, req entity.OrderEntity) (int64, error) {
	orderDate, orderTime, err := conv.ParseStringToDateTime(req.OrderDate, req.OrderTime)
	if err != nil {
		o.logger.Errorf("[OrderRepository] CreateOrder: %v", err)
		return 0, err
	}

	var (
		db         = o.getDB(ctx)
		modelOrder = &model.Order{
			OrderCode:    req.OrderCode,
			BuyerID:      req.BuyerID,
			OrderDate:    *orderDate,
			OrderTime:    *orderTime,
			Status:       req.Status,
			TotalAmount:  float64(req.TotalAmount),
			ShippingType: req.ShippingType,
			ShippingFee:  float64(req.ShippingFee),
			Remarks:      req.Remarks,
			BuyerName:    req.BuyerName,
			BuyerEmail:   req.BuyerEmail,
			BuyerPhone:   req.BuyerPhone,
			BuyerAddress: req.BuyerAddress,
			BuyerLat:     req.BuyerLat,
			BuyerLng:     req.BuyerLng,
		}
	)

	for _, item := range req.OrderItems {
		modelOrder.OrderItems = append(modelOrder.OrderItems, model.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	for _, snapshot := range req.ProductSnapshots {
		modelOrder.ProductSnapshots = append(modelOrder.ProductSnapshots, model.ProductSnapshot{
			ProductID:    snapshot.ProductID,
			Name:         snapshot.Name,
			Image:        snapshot.Image,
			RegularPrice: snapshot.RegularPrice,
			SalePrice:    snapshot.SalePrice,
			Unit:         snapshot.Unit,
			Weight:       snapshot.Weight,
		})
	}

	if err := db.WithContext(ctx).Create(&modelOrder).Error; err != nil {
		o.logger.Errorf("[OrderRepository] CreateOrder: %v", err)
		if strings.Contains(err.Error(), "foreign key") {
			err = utils.ErrRelationDataNotFound
		}
		return 0, err
	}

	return modelOrder.ID, nil
}

// DeleteOrder implements [OrderRepositoryInterface].
func (o *orderRepository) DeleteOrder(ctx context.Context, orderId int64) error {
	panic("unimplemented")
}

// GetAllOrders implements [OrderRepositoryInterface].
func (o *orderRepository) GetAllOrders(ctx context.Context, query entity.OrderQueryString) ([]entity.OrderEntity, int64, int64, error) {
	var (
		db          = o.getDB(ctx)
		modelOrders []model.Order
		entities    []entity.OrderEntity
		countData   int64
	)

	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_id", "quantity")
		}).
		Preload("ProductSnapshots", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "last_used", "regular_price")
		})

	if query.Search != "" {
		sqlMain = sqlMain.Where(`order_code ILIKE ?`, "%"+query.Search+"%")
	}

	if query.Status != "" {
		sqlMain = sqlMain.Where("status = ?", query.Status)
	}

	if query.BuyerID != 0 {
		sqlMain = sqlMain.Where("buyer_id = ?", query.BuyerID)
	}

	if err := sqlMain.Model(&modelOrders).Count(&countData).Error; err != nil {
		o.logger.Errorf("[OrderRepository] GetAllOrders: %v", err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order("id DESC").
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelOrders).Error; err != nil {
		o.logger.Errorf("[OrderRepository] GetAllOrders: %v", err)
		return nil, 0, 0, err
	}

	for _, val := range modelOrders {
		snapshotMap := make(map[int64]entity.ProductSnapshotEntity)

		orderEntity := entity.OrderEntity{
			ID:           val.ID,
			OrderCode:    val.OrderCode,
			Status:       val.Status,
			OrderDate:    val.OrderDate.Format("2006-01-02"),
			OrderTime:    val.OrderTime.Format("15:04:05"),
			TotalAmount:  int64(val.TotalAmount),
			BuyerID:      val.BuyerID,
			BuyerName:    val.BuyerName,
			BuyerEmail:   val.BuyerEmail,
			BuyerPhone:   val.BuyerPhone,
			BuyerAddress: val.BuyerAddress,
			BuyerLat:     val.BuyerLat,
			BuyerLng:     val.BuyerLng,
		}

		for _, snapshot := range val.ProductSnapshots {
			snapshotMap[snapshot.ProductID] = entity.ProductSnapshotEntity{
				ProductID: snapshot.ProductID,
				Name:      snapshot.Name,
				Image:     snapshot.Image,
				SalePrice: snapshot.SalePrice,
				Weight:    snapshot.Weight,
				Unit:      snapshot.Unit,
			}
		}

		for _, item := range val.OrderItems {
			if snapshot, ok := snapshotMap[item.ProductID]; ok {
				orderEntity.OrderItems = append(orderEntity.OrderItems, entity.OrderItemEntity{
					ID:        item.ID,
					OrderID:   item.OrderID,
					ProductID: item.ProductID,
					Quantity:  item.Quantity,

					ProductName:  snapshot.Name,
					ProductImage: snapshot.Image,
					Price:        int64(snapshot.SalePrice),
				})
			}
		}

		entities = append(entities, orderEntity)
	}

	return entities, countData, int64(totalPages), nil
}

// GetOrderById implements [OrderRepositoryInterface].
func (o *orderRepository) GetOrderById(ctx context.Context, orderId int64, userId int64) (*entity.OrderEntity, error) {
	var (
		db          = o.getDB(ctx)
		modelOrder  model.Order
		orderEntity *entity.OrderEntity
		snapshotMap = make(map[int64]entity.ProductSnapshotEntity)
	)

	sqlMain := db.WithContext(ctx).
		Where("id = ?", orderId).
		Omit("created_at", "updated_at", "deleted_at").
		Preload("OrderItems", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "order_id", "product_id", "quantity")
		}).
		Preload("ProductSnapshots", func(db *gorm.DB) *gorm.DB {
			return db.Omit("created_at", "last_used", "regular_price")
		})

	if err := sqlMain.First(&modelOrder).Error; err != nil {
		o.logger.Errorf("[OrderRepository] GetOrderById: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = utils.ErrDataNotFound
		}
		return nil, err
	}

	if userId != 0 && modelOrder.BuyerID != userId {
		err := utils.ErrAccessForbidden
		o.logger.Errorf("[OrderRepository] GetOrderById: %v", err)
		return nil, err
	}

	orderEntity = &entity.OrderEntity{
		ID:           modelOrder.ID,
		OrderCode:    modelOrder.OrderCode,
		OrderDate:    modelOrder.OrderDate.Format("2006-01-02"),
		OrderTime:    modelOrder.OrderTime.Format("15:04:05"),
		Status:       modelOrder.Status,
		BuyerID:      modelOrder.BuyerID,
		TotalAmount:  int64(modelOrder.TotalAmount),
		Remarks:      modelOrder.Remarks,
		ShippingType: modelOrder.ShippingType,
		ShippingFee:  int64(modelOrder.ShippingFee),
		BuyerName:    modelOrder.BuyerName,
		BuyerEmail:   modelOrder.BuyerEmail,
		BuyerPhone:   modelOrder.BuyerPhone,
		BuyerAddress: modelOrder.BuyerAddress,
		BuyerLat:     modelOrder.BuyerLat,
		BuyerLng:     modelOrder.BuyerLng,
	}

	for _, snapshot := range modelOrder.ProductSnapshots {
		snapshotMap[snapshot.ProductID] = entity.ProductSnapshotEntity{
			ProductID: snapshot.ProductID,
			Name:      snapshot.Name,
			Image:     snapshot.Image,
			SalePrice: snapshot.SalePrice,
			Weight:    snapshot.Weight,
			Unit:      snapshot.Unit,
		}
	}

	for _, item := range modelOrder.OrderItems {
		if snapshot, ok := snapshotMap[item.ProductID]; ok {
			orderEntity.OrderItems = append(orderEntity.OrderItems, entity.OrderItemEntity{
				ID:        item.ID,
				OrderID:   item.OrderID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,

				ProductName:  snapshot.Name,
				ProductImage: snapshot.Image,
				Price:        int64(snapshot.SalePrice),
			})
		}
	}

	return orderEntity, nil
}

// UpdateOrderStatus implements [OrderRepositoryInterface].
func (o *orderRepository) UpdateOrderStatus(ctx context.Context, req entity.OrderEntity) error {
	var (
		db         = o.getDB(ctx)
		modelOrder = model.Order{
			ID:      req.ID,
			Status:  req.Status,
			Remarks: req.Remarks,
		}
	)

	if err := db.WithContext(ctx).Updates(&modelOrder).Error; err != nil {
		o.logger.Errorf("[OrderRepository] UpdateOrderStatus: %v", err)
		return err
	}

	return nil
}
