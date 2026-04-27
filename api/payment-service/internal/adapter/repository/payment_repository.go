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
	"gorm.io/gorm/clause"
)

type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *entity.PaymentEntity) (uint, string, error)
	CreatePaymentLog(ctx context.Context, paymentId uint, status string) error
	GetPaymentById(ctx context.Context, paymentId uint, userId uint) (*entity.PaymentEntity, error)
	UpdateStatusByOrderCode(ctx context.Context, orderId uint, status string) (uint, error)
	GetAllPayments(ctx context.Context, query entity.QueryStringPayment) ([]entity.PaymentEntity, int64, int64, error)
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
func (p *paymentRepository) UpdateStatusByOrderCode(ctx context.Context, orderId uint, status string) (uint, error) {
	var (
		db           = p.getDB(ctx)
		modelPayment = model.Payment{
			PaymentStatus: status,
		}
	)

	tx := db.WithContext(ctx).
		Model(&modelPayment).
		Clauses(clause.Returning{
			Columns: []clause.Column{
				{Name: "id"},
			},
		}).
		Where("order_id = ?", orderId).
		Updates(&modelPayment)

	if tx.Error != nil {
		p.logger.Errorf("[PaymentRepository-1] UpdateStatusByOrderCode: %v", tx.Error)
		return 0, tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		p.logger.Errorf("[PaymentRepository-2] UpdateStatusByOrderCode: %v", err)
		return 0, err
	}

	return modelPayment.ID, nil
}

// GetPaymentById implements [PaymentRepositoryInterface].
func (p *paymentRepository) GetPaymentById(ctx context.Context, paymentId uint, userId uint) (*entity.PaymentEntity, error) {
	var (
		db            = p.getDB(ctx)
		modelPayment  model.Payment
		paymentEntity entity.PaymentEntity
	)

	sqlMain := db.Debug().WithContext(ctx).
		Omit("updated_at", "deleted_at").
		Preload("PaymentLogs", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "payment_id", "status")
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
		OrderID:          modelPayment.OrderID,
		UserID:           modelPayment.UserID,
		PaymentMethod:    modelPayment.PaymentMethod,
		PaymentStatus:    modelPayment.PaymentStatus,
		PaymentGatewayID: *modelPayment.PaymentGatewayID,
		GrossAmount:      modelPayment.GrossAmount,
		PaymentURL:       *modelPayment.PaymentURL,
		PaymentAt:        modelPayment.CreatedAt.Format("2006-01-02 15:05:05"),
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
func (p *paymentRepository) CreatePaymentLog(ctx context.Context, paymentId uint, status string) error {
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
func (p *paymentRepository) CreatePayment(ctx context.Context, payment *entity.PaymentEntity) (uint, string, error) {
	var (
		db           = p.getDB(ctx)
		modelPayment = model.Payment{
			OrderID:          payment.OrderID,
			UserID:           payment.UserID,
			PaymentMethod:    payment.PaymentMethod,
			PaymentStatus:    payment.PaymentStatus,
			PaymentGatewayID: &payment.PaymentGatewayID,
			GrossAmount:      payment.GrossAmount,
			PaymentURL:       &payment.PaymentURL,
		}
	)

	if err := db.WithContext(ctx).Create(&modelPayment).Error; err != nil {
		p.logger.Errorf("[PaymentRepository-1] CreatePayment: %v", err)
		return 0, "", err
	}

	return modelPayment.ID, modelPayment.PaymentStatus, nil
}
