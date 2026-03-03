package repository

import (
	"context"
	"errors"
	"time"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
	"user-service/utils"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type VerificationTokenRepositoryInterface interface {
	CreateVerificationToken(ctx context.Context, req entity.VerificationTokenEntity) error
	GetDataByToken(ctx context.Context, token string) (*entity.VerificationTokenEntity, error)
	DeleteVerificationToken(ctx context.Context, tokenId int64) error

	getDB(ctx context.Context) *gorm.DB
}

type verificationTokenRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewVerificationTokenRepository(db *gorm.DB, logger *log.Logger) VerificationTokenRepositoryInterface {
	return &verificationTokenRepository{db: db, logger: logger}
}

// getDB implements [VerificationTokenRepositoryInterface].
func (v *verificationTokenRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return v.db
}

// DeleteVerificationToken implements VerificationTokenRepositoryInterface.
func (v *verificationTokenRepository) DeleteVerificationToken(ctx context.Context, tokenId int64) error {
	var (
		db                     = v.getDB(ctx)
		modelVerificationToken = model.VerificationToken{
			ID: tokenId,
		}
	)

	if err := db.WithContext(ctx).Delete(&modelVerificationToken).Error; err != nil {
		log.Errorf("[VerificationTokenRepository-3] DeleteVerificationToken: %v", err)
		return err
	}

	return nil
}

// GetDataByToken implements VerificationTokenRepositoryInterface.
func (v *verificationTokenRepository) GetDataByToken(ctx context.Context, token string) (*entity.VerificationTokenEntity, error) {
	var (
		db       = v.getDB(ctx)
		tokenDTO model.VerificationTokenDTO
	)

	tokenFieldSelect := `
		verification_tokens.id AS id,
		verification_tokens.user_id AS user_id,
		verification_tokens.token_type AS token_type,
		verification_tokens.expires_at AS expires_at,
		users.name AS user_name,
		users.email AS user_email,
		roles.name AS role_name
	`
	if err := db.WithContext(ctx).
		Model(&model.VerificationToken{}).
		Select(tokenFieldSelect).
		Joins("INNER JOIN users ON verification_tokens.user_id = users.id").
		Joins("INNER JOIN user_role ON users.id = user_role.user_id").
		Joins("INNER JOIN roles ON user_role.role_id = roles.id").
		Where("verification_tokens.token = ?", token).
		Where("roles.name = ?", "Customer").
		First(&tokenDTO).Error; err != nil {
		log.Errorf("[VerificationTokenRepository-1] GetDataByToken: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.TOKEN_INVALID)
			return nil, err
		}
		return nil, err
	}

	return &entity.VerificationTokenEntity{
		ID:        tokenDTO.ID,
		UserID:    tokenDTO.UserID,
		TokenType: tokenDTO.TokenType,
		ExpiresAt: tokenDTO.ExpiresAt,
		User: entity.UserEntity{
			ID:       tokenDTO.UserID,
			Email:    tokenDTO.UserEmail,
			Name:     tokenDTO.UserName,
			RoleName: tokenDTO.RoleName,
		},
	}, nil
}

// CreateVerificationToken implements VerificationTokenRepositoryInterface.
func (v *verificationTokenRepository) CreateVerificationToken(ctx context.Context, req entity.VerificationTokenEntity) error {
	var (
		db                     = v.getDB(ctx)
		modelVerificationToken = model.VerificationToken{
			UserID:    req.UserID,
			Token:     req.Token,
			TokenType: req.TokenType,
			ExpiresAt: time.Now().Add(1 * time.Hour),
		}
	)

	if err := db.WithContext(ctx).Create(&modelVerificationToken).Error; err != nil {
		log.Errorf("[VerificationTokenRepository-1] CreateVerificationToken: %v", err)
		return err
	}

	return nil
}
