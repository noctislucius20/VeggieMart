package repository

import (
	"context"
	"errors"
	"math"
	"payment-service/internal/core/domain/entity"
	"payment-service/internal/core/domain/model"
	"payment-service/utils"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *entity.PaymentEntity, db *gorm.DB) (uint, string, error)
	CreatePaymentLog(ctx context.Context, paymentId uint, status string, db *gorm.DB) error
	GetPaymentById(ctx context.Context, paymentId uint, userId uint, db *gorm.DB) (*entity.PaymentEntity, error)
	UpdateStatusByOrderCode(ctx context.Context, orderId uint, status string, db *gorm.DB) error
	GetAllPayments(ctx context.Context, query entity.QueryStringPayment, db *gorm.DB) ([]entity.PaymentEntity, int64, int64, error)
}

type paymentRepository struct {
	logger *log.Logger
}

// GetAllPayments implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetAllPayments(ctx context.Context, query entity.QueryStringPayment, db *gorm.DB) ([]entity.PaymentEntity, int64, int64, error) {
	modelPayments := []model.Payment{}

	var countData int64

	offset := (query.Page - 1) * query.Limit

	sqlMain := db.WithContext(ctx).
		Where("payment_method ILIKE ? OR payment_status ILIKE ?", "%"+query.Search+"%", "%"+query.Status+"%")

	if query.UserID != 0 {
		sqlMain = sqlMain.Where("user_id = ?", query.UserID)
	}

	if err := sqlMain.Model(&modelPayments).Count(&countData).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] GetAllPayments: %v", err.Error())
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))
	if err := sqlMain.Order("id DESC").
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelPayments).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-2] GetAllPayments: %v", err.Error())
		return nil, 0, 0, err
	}

	entities := []entity.PaymentEntity{}
	for _, val := range modelPayments {
		entities = append(entities, entity.PaymentEntity{
			ID:               val.ID,
			OrderID:          val.OrderID,
			UserID:           val.UserID,
			PaymentMethod:    val.PaymentMethod,
			PaymentStatus:    val.PaymentStatus,
			PaymentGatewayID: *val.PaymentGatewayID,
			GrossAmount:      val.GrossAmount,
			PaymentURL:       *val.PaymentURL,
		})
	}

	return entities, countData, int64(totalPages), nil
}

// UpdateStatusByOrderCode implements [PaymentRepositoryInterface].
func (p *paymentRepository) UpdateStatusByOrderCode(ctx context.Context, orderId uint, status string, db *gorm.DB) error {
	modelPayment := model.Payment{
		PaymentStatus: status,
	}

	if err := db.WithContext(ctx).
		Where("order_id = ?", orderId).
		Updates(&modelPayment).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] UpdateStatusByOrderCode: %v", err.Error())
		return err
	}

	return nil
}

// GetPaymentById implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetPaymentById(ctx context.Context, paymentId uint, userId uint, db *gorm.DB) (*entity.PaymentEntity, error) {
	modelPayment := model.Payment{}

	sqlMain := db.WithContext(ctx).
		Omit("updated_at", "deleted_at").
		Preload("PaymentLogs", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "status")
		})

	if userId != 0 {
		sqlMain = sqlMain.Where("user_id", userId)
	}

	if err := sqlMain.
		First(&modelPayment, "id = ?", paymentId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New(utils.DATA_NOT_FOUND)
		}
		p.logger.Errorf("[PaymentRepository-1] GetPaymentById: %v", err.Error())
		return nil, err
	}

	paymentLogEntities := []entity.PaymentLogEntity{}
	for _, item := range modelPayment.PaymentLogs {
		paymentLogEntities = append(paymentLogEntities, entity.PaymentLogEntity{
			ID:        item.ID,
			PaymentID: item.PaymentID,
			Status:    item.Status,
		})
	}

	return &entity.PaymentEntity{
		ID:               modelPayment.ID,
		OrderID:          modelPayment.OrderID,
		UserID:           modelPayment.UserID,
		PaymentMethod:    modelPayment.PaymentMethod,
		PaymentStatus:    modelPayment.PaymentStatus,
		PaymentGatewayID: *modelPayment.PaymentGatewayID,
		GrossAmount:      modelPayment.GrossAmount,
		PaymentURL:       *modelPayment.PaymentURL,
		PaymentAt:        modelPayment.CreatedAt.Format("2006-01-02 15:05:05"),
		PaymentLogs:      paymentLogEntities,
	}, nil
}

// CreatePaymentLog implements [PaymentRepositoryInterface].
func (p *paymentRepository) CreatePaymentLog(ctx context.Context, paymentId uint, status string, db *gorm.DB) error {
	modelPaymentLog := model.PaymentLog{
		PaymentID: paymentId,
		Status:    status,
	}

	if err := db.WithContext(ctx).Create(&modelPaymentLog).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] CreatePaymentLog: %v", err)
		return err
	}

	return nil
}

// CreatePayment implements [PaymentRepositoryInterface].
func (p *paymentRepository) CreatePayment(ctx context.Context, payment *entity.PaymentEntity, db *gorm.DB) (uint, string, error) {
	modelPayment := model.Payment{
		OrderID:          payment.OrderID,
		UserID:           payment.UserID,
		PaymentMethod:    payment.PaymentMethod,
		PaymentStatus:    payment.PaymentStatus,
		PaymentGatewayID: &payment.PaymentGatewayID,
		GrossAmount:      payment.GrossAmount,
		PaymentURL:       &payment.PaymentURL,
	}

	if err := db.WithContext(ctx).Create(&modelPayment).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] CreatePayment: %v", err)
		return 0, "", err
	}

	return modelPayment.ID, modelPayment.PaymentStatus, nil
}

func NewPaymentRepository(logger *log.Logger) PaymentRepositoryInterface {
	return &paymentRepository{logger: logger}
}
