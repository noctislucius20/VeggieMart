package repository

import (
	"context"
	"errors"
	"strings"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
	"user-service/utils"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetRolesAll(ctx context.Context, search string) ([]entity.RoleEntity, error)
	GetRoleByIdOrName(ctx context.Context, id int64, name string) (*entity.RoleEntity, error)
	CreateRole(ctx context.Context, req entity.RoleEntity) (int64, error)
	DeleteRole(ctx context.Context, id int64) error
	UpdateRole(ctx context.Context, req entity.RoleEntity) error

	getDB(ctx context.Context) *gorm.DB
}

type roleRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewRoleRepository(db *gorm.DB, logger *log.Logger) RoleRepositoryInterface {
	return &roleRepository{
		logger: logger,
	}
}

// getDB implements [RoleRepositoryInterface].
func (r *roleRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return r.db
}

// CreateRole implements RoleRepositoryInterface.
func (r *roleRepository) CreateRole(ctx context.Context, req entity.RoleEntity) (int64, error) {
	var (
		db        = r.getDB(ctx)
		modelRole = model.Role{
			Name: req.Name,
		}
	)

	if err := db.WithContext(ctx).Create(&modelRole).Error; err != nil {
		r.logger.Errorf("[RoleRepository-1] CreateRole: %v", err)
		if strings.Contains(err.Error(), "duplicate key") {
			err := errors.New(utils.DATA_ALREADY_EXISTS)
			return 0, err
		}
		return 0, err
	}

	return modelRole.ID, nil
}

// DeleteRole implements RoleRepositoryInterface.
func (r *roleRepository) DeleteRole(ctx context.Context, id int64) error {
	var (
		db        = r.getDB(ctx)
		modelRole = model.Role{
			ID: id,
		}
		roleDeleteDTO model.RoleDeleteDTO
	)

	if err := db.WithContext(ctx).
		Model(&model.Role{}).
		Select("roles.id AS roles_id",
			"user_role.role_id AS user_role_role_id",
		).
		Joins("LEFT JOIN user_role ON roles.id = user_role.role_id").
		Where("roles.id = ?", id).
		First(&roleDeleteDTO).Error; err != nil {
		r.logger.Errorf("[RoleRepository-1] DeleteRole: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return err
		}
		return err
	}

	if roleDeleteDTO.UserRoleRoleID > 0 {
		err := errors.New(utils.DATA_STILL_IN_USED)
		r.logger.Errorf("[RoleRepository-2] DeleteRole: %v", err)
		return err
	}

	if err := db.WithContext(ctx).Delete(&modelRole).Error; err != nil {
		r.logger.Errorf("[RoleRepository-3] DeleteRole: %v", err)
		return err
	}

	return nil
}

// GetRolesAll implements RoleRepositoryInterface.
func (r *roleRepository) GetRolesAll(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	var (
		db         = r.getDB(ctx)
		modelRoles []model.Role
	)

	if err := db.WithContext(ctx).Select("id", "name").Where("name ILIKE ?", "%"+search+"%").Find(&modelRoles).Error; err != nil {
		r.logger.Errorf("[RoleRepository-1] GetRolesAll: %v", err)
		return nil, err
	}

	entityRole := []entity.RoleEntity{}
	for _, modelRole := range modelRoles {
		entityRole = append(entityRole, entity.RoleEntity{
			ID:   modelRole.ID,
			Name: modelRole.Name,
		})
	}

	return entityRole, nil
}

// GetById implements RoleRepositoryInterface.
func (r *roleRepository) GetRoleByIdOrName(ctx context.Context, id int64, name string) (*entity.RoleEntity, error) {
	var (
		db        = r.getDB(ctx)
		modelRole model.Role
	)

	if err := db.WithContext(ctx).Select("id", "name").
		First(&modelRole, "id = ? OR name = ?", id, name).Error; err != nil {
		r.logger.Errorf("[RoleRepository-1] GetRoleByIdOrName: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		}
		return nil, err
	}

	return &entity.RoleEntity{
		ID:   modelRole.ID,
		Name: modelRole.Name,
	}, nil
}

// UpdateRole implements RoleRepositoryInterface.
func (r *roleRepository) UpdateRole(ctx context.Context, req entity.RoleEntity) error {
	var (
		db        = r.getDB(ctx)
		modelRole = model.Role{
			ID:   req.ID,
			Name: req.Name,
		}
	)

	tx := db.WithContext(ctx).Updates(&modelRole)
	if tx.Error != nil {
		r.logger.Errorf("[RoleRepository-1] UpdateRole: %v", tx.Error)
		if strings.Contains(tx.Error.Error(), "duplicate key") {
			err := errors.New(utils.DATA_ALREADY_EXISTS)
			return err
		}
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		r.logger.Errorf("[RoleRepository-2] UpdateRole: %v", err)
		return err
	}

	return nil
}
