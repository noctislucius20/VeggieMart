package repository

import (
	"context"
	"errors"
	"math"
	"payment-service/internal/core/domain/entity"
	"payment-service/internal/core/domain/model"
	"payment-service/utils"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *entity.PaymentEntity) (int64, string, error)
	CreatePaymentLog(ctx context.Context, paymentId int64, status string) error
	GetPaymentById(ctx context.Context, paymentId int64, userId int64) (*entity.PaymentEntity, error)
	GetPaymentByOrderId(ctx context.Context, orderId int64, userId int64) (*entity.PaymentEntity, error)
	UpdateStatusByPaymentId(ctx context.Context, paymentId int64, status string) error
	UpdatePaymentGatewayByOrderId(ctx context.Context, orderId int64, paymentGatewayId *string, paymentUrl *string) (int64, error)
	GetAllPayments(ctx context.Context, query entity.QueryStringPayment) ([]entity.PaymentEntity, int64, int64, error)
	GetPaymentOrderIdByOrderCode(ctx context.Context, orderCode string) (int64, int64, error)
}

type paymentRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewPaymentRepository(db *gorm.DB, logger *log.Logger) PaymentRepositoryInterface {
	return &paymentRepository{db: db, logger: logger}
}

// getDB implements [PaymentRepositoryInterface].
func (p *paymentRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return p.db
}

// UpdatePaymentGatewayByOrderId implements [PaymentRepositoryInterface].
func (p *paymentRepository) UpdatePaymentGatewayByOrderId(ctx context.Context, orderId int64, paymentGatewayId *string, paymentUrl *string) (int64, error) {
	var (
		db           = p.getDB(ctx)
		modelPayment model.Payment
	)

	if err := db.WithContext(ctx).
		Select("payments.id").
		InnerJoins("OrderSnapshot", db.
			Select("id", "payment_id", "order_id")).
		Where(`"OrderSnapshot"."order_id" = ?`, orderId).
		Where("payments.payment_status = ?", "PENDING").
		First(&modelPayment).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] UpdatePaymentGatewayByOrderId: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		return 0, err
	}

	modelPayment.PaymentGatewayID = paymentGatewayId
	modelPayment.PaymentURL = paymentUrl

	if err := db.WithContext(ctx).Updates(&modelPayment).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-2] UpdatePaymentGatewayByOrderId: %v", err)
		return 0, err
	}

	return modelPayment.ID, nil
}

// GetPaymentOrderIdByOrderCode implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetPaymentOrderIdByOrderCode(ctx context.Context, orderCode string) (int64, int64, error) {
	var (
		db                 = p.getDB(ctx)
		modelOrderSnapshot model.OrderSnapshot
	)

	if err := db.WithContext(ctx).
		Select("id", "payment_id", "order_id").
		Where("order_code = ?", orderCode).
		Order("id DESC").
		First(&modelOrderSnapshot).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] GetPaymentOrderIdByOrderCode: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		return 0, 0, err
	}

	return modelOrderSnapshot.PaymentID, modelOrderSnapshot.OrderID, nil
}

// GetAllPayments implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetAllPayments(ctx context.Context, query entity.QueryStringPayment) ([]entity.PaymentEntity, int64, int64, error) {
	var (
		db            = p.getDB(ctx)
		modelPayments []model.Payment
		entities      []entity.PaymentEntity
		countData     int64
	)

	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).
		Preload("OrderSnapshot", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "order_id", "shipping_type", "order_code")
		}).
		Where("payment_method ILIKE ? OR payment_status ILIKE ?", "%"+query.Search+"%", "%"+query.Status+"%")

	if query.UserID != 0 {
		sqlMain = sqlMain.Where("user_id = ?", query.UserID)
	}

	if err := sqlMain.Model(&modelPayments).Count(&countData).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] GetAllPayments: %v", err)
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order("id DESC").
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelPayments).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-2] GetAllPayments: %v", err)
		return nil, 0, 0, err
	}

	for _, val := range modelPayments {
		entities = append(entities, entity.PaymentEntity{
			ID:               val.ID,
			UserID:           val.UserID,
			PaymentMethod:    val.PaymentMethod,
			PaymentStatus:    val.PaymentStatus,
			PaymentGatewayID: *val.PaymentGatewayID,
			GrossAmount:      val.GrossAmount,
			PaymentURL:       *val.PaymentURL,
			Order: entity.OrderEntity{
				ID:           val.OrderSnapshot.OrderID,
				OrderCode:    val.OrderSnapshot.OrderCode,
				ShippingType: val.OrderSnapshot.ShippingType,
			},
		})
	}

	return entities, countData, int64(totalPages), nil
}

// UpdateStatusByPaymentId implements [PaymentRepositoryInterface].
func (p *paymentRepository) UpdateStatusByPaymentId(ctx context.Context, paymentId int64, status string) error {
	var (
		db           = p.getDB(ctx)
		modelPayment = model.Payment{
			ID:            paymentId,
			PaymentStatus: status,
		}
	)

	tx := db.WithContext(ctx).Updates(&modelPayment)

	if tx.Error != nil {
		p.logger.Errorf("[PaymentRepository-1] UpdateStatusByPaymentId: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		p.logger.Errorf("[PaymentRepository-2] UpdateStatusByPaymentId: %v", err)
		return err
	}

	return nil
}

// GetPaymentById implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetPaymentById(ctx context.Context, paymentId int64, userId int64) (*entity.PaymentEntity, error) {
	var (
		db            = p.getDB(ctx)
		modelPayment  model.Payment
		paymentEntity entity.PaymentEntity
	)

	sqlMain := db.WithContext(ctx).
		Omit("updated_at", "deleted_at").
		Preload("PaymentLogs", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "status")
		}).
		Preload("OrderSnapshot", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "order_id", "order_code", "shipping_type", "remarks", "order_datetime", "customer_name", "customer_address", "customer_email")
		})

	if userId != 0 {
		sqlMain = sqlMain.Where("user_id = ?", userId)
	}

	if err := sqlMain.
		First(&modelPayment, "id = ?", paymentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		p.logger.Errorf("[PaymentRepository-1] GetPaymentById: %v", err)
		return nil, err
	}

	paymentEntity = entity.PaymentEntity{
		ID:               modelPayment.ID,
		UserID:           modelPayment.UserID,
		PaymentMethod:    modelPayment.PaymentMethod,
		PaymentStatus:    modelPayment.PaymentStatus,
		PaymentGatewayID: *modelPayment.PaymentGatewayID,
		GrossAmount:      modelPayment.GrossAmount,
		PaymentURL:       *modelPayment.PaymentURL,
		PaymentAt:        modelPayment.CreatedAt.Format("2006-01-02 15:05:05"),
		Order: entity.OrderEntity{
			ID:            modelPayment.OrderSnapshot.OrderID,
			OrderCode:     modelPayment.OrderSnapshot.OrderCode,
			ShippingType:  modelPayment.OrderSnapshot.ShippingType,
			OrderDatetime: modelPayment.OrderSnapshot.OrderDatetime.Format("2006-01-02 15:05:05"),
			Remarks:       modelPayment.OrderSnapshot.Remarks,
		},
		Customer: entity.CustomerEntity{
			CustomerName:    modelPayment.OrderSnapshot.CustomerName,
			CustomerEmail:   modelPayment.OrderSnapshot.CustomerEmail,
			CustomerAddress: modelPayment.OrderSnapshot.CustomerAddress,
		},
	}

	for _, item := range modelPayment.PaymentLogs {
		paymentEntity.PaymentLogs = append(paymentEntity.PaymentLogs, entity.PaymentLogEntity{
			ID:        item.ID,
			PaymentID: item.PaymentID,
			Status:    item.Status,
		})
	}

	return &paymentEntity, nil
}

// GetPaymentByOrderId implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetPaymentByOrderId(ctx context.Context, orderId int64, userId int64) (*entity.PaymentEntity, error) {
	var (
		db            = p.getDB(ctx)
		modelPayment  model.Payment
		paymentEntity entity.PaymentEntity
	)

	sqlMain := db.WithContext(ctx).
		Omit("updated_at", "deleted_at").
		InnerJoins("OrderSnapshot", db.Select("id", "payment_id", "order_id", "order_code", "shipping_type", "remarks", "order_datetime", "customer_name", "customer_address", "customer_email").
			Where("order_id = ?", orderId).
			Model(&model.OrderSnapshot{})).
		Preload("PaymentLogs", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "status")
		})

	if userId != 0 {
		sqlMain = sqlMain.Where("user_id = ?", userId)
	}

	if err := sqlMain.
		First(&modelPayment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		p.logger.Errorf("[PaymentRepository-1] GetPaymentByOrderId: %v", err)
		return nil, err
	}

	paymentEntity = entity.PaymentEntity{
		ID:               modelPayment.ID,
		UserID:           modelPayment.UserID,
		PaymentMethod:    modelPayment.PaymentMethod,
		PaymentStatus:    modelPayment.PaymentStatus,
		PaymentGatewayID: *modelPayment.PaymentGatewayID,
		GrossAmount:      modelPayment.GrossAmount,
		PaymentURL:       *modelPayment.PaymentURL,
		PaymentAt:        modelPayment.CreatedAt.Format("2006-01-02 15:05:05"),
		Order: entity.OrderEntity{
			ID:            modelPayment.OrderSnapshot.OrderID,
			OrderCode:     modelPayment.OrderSnapshot.OrderCode,
			ShippingType:  modelPayment.OrderSnapshot.ShippingType,
			OrderDatetime: modelPayment.OrderSnapshot.OrderDatetime.Format("2006-01-02 15:05:05"),
			Remarks:       modelPayment.OrderSnapshot.Remarks,
		},
		Customer: entity.CustomerEntity{
			CustomerName:    modelPayment.OrderSnapshot.CustomerName,
			CustomerEmail:   modelPayment.OrderSnapshot.CustomerEmail,
			CustomerAddress: modelPayment.OrderSnapshot.CustomerAddress,
		},
	}

	for _, item := range modelPayment.PaymentLogs {
		paymentEntity.PaymentLogs = append(paymentEntity.PaymentLogs, entity.PaymentLogEntity{
			ID:        item.ID,
			PaymentID: item.PaymentID,
			Status:    item.Status,
		})
	}

	return &paymentEntity, nil
}

// CreatePaymentLog implements [PaymentRepositoryInterface].
func (p *paymentRepository) CreatePaymentLog(ctx context.Context, paymentId int64, status string) error {
	var (
		db              = p.getDB(ctx)
		modelPaymentLog = model.PaymentLog{
			PaymentID: paymentId,
			Status:    status,
		}
	)

	if err := db.WithContext(ctx).Create(&modelPaymentLog).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] CreatePaymentLog: %v", err)
		return err
	}

	return nil
}

// CreatePayment implements [PaymentRepositoryInterface].
func (p *paymentRepository) CreatePayment(ctx context.Context, payment *entity.PaymentEntity) (int64, string, error) {
	parsedDatetime, err := time.Parse("2006-01-02 15:04:05", payment.Order.OrderDate+" "+payment.Order.OrderTime)
	if err != nil {
		return 0, "", err
	}

	var (
		db           = p.getDB(ctx)
		modelPayment = model.Payment{
			UserID:        int64(payment.Order.BuyerID),
			PaymentMethod: payment.Order.PaymentMethod,
			PaymentStatus: payment.PaymentStatus,
			GrossAmount:   float64(payment.Order.TotalAmount),
			OrderSnapshot: &model.OrderSnapshot{
				OrderID:         payment.Order.ID,
				OrderCode:       payment.Order.OrderCode,
				OrderDatetime:   parsedDatetime,
				Status:          payment.Order.Status,
				PaymentMethod:   payment.Order.PaymentMethod,
				ShippingFee:     float64(payment.Order.ShippingFee),
				ShippingType:    payment.Order.ShippingType,
				Remarks:         payment.Order.Remarks,
				TotalAmount:     float64(payment.Order.TotalAmount),
				CustomerID:      payment.Order.BuyerID,
				CustomerName:    payment.Order.BuyerName,
				CustomerEmail:   payment.Order.BuyerEmail,
				CustomerAddress: payment.Order.BuyerAddress,
				CustomerPhone:   payment.Order.BuyerPhone,
			},
		}
	)

	if err := db.WithContext(ctx).Create(&modelPayment).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] CreatePayment: %v", err)
		return 0, "", err
	}

	return modelPayment.ID, modelPayment.PaymentStatus, nil
}
