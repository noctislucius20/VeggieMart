package service

import (
	"context"
	"user-service/internal/adapter/repository"
	"user-service/internal/adapter/repository/cache"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service/transaction"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type RoleServiceInterface interface {
	GetRolesAllAdmin(ctx context.Context, search string) ([]entity.RoleEntity, error)
	GetRoleByIdAdmin(ctx context.Context, id int64) (*entity.RoleEntity, error)
	GetRoleByNameAdmin(ctx context.Context, name string) (*entity.RoleEntity, error)
	CreateRoleAdmin(ctx context.Context, req entity.RoleEntity) (int64, error)
	DeleteRoleAdmin(ctx context.Context, id int64) error
	UpdateRoleAdmin(ctx context.Context, req entity.RoleEntity) error
}

type roleService struct {
	repo        repository.RoleRepositoryInterface
	redisClient *redis.Client
	cacheRole   cache.RoleCacheInterface
	txManager   transaction.TransactionManager
	logger      *log.Logger
}

func NewRoleService(repo repository.RoleRepositoryInterface, redisClient *redis.Client, cacheRole cache.RoleCacheInterface, txManager transaction.TransactionManager, logger *log.Logger) RoleServiceInterface {
	return &roleService{
		repo:        repo,
		redisClient: redisClient,
		cacheRole:   cacheRole,
		txManager:   txManager,
		logger:      logger,
	}
}

// GetRoleByNameAdmin implements [RoleServiceInterface].
func (r *roleService) GetRoleByNameAdmin(ctx context.Context, name string) (*entity.RoleEntity, error) {
	var (
		role *entity.RoleEntity
	)

	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := r.cacheRole.GetRoleByName(txCtx, name)
		if err != nil {
			return err
		}

		role = roleEntity

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] GetRoleByNameAdmin: %v", err)
		return nil, err
	}

	return role, nil
}

// CreateRoleAdmin implements RoleServiceInterface.
func (r *roleService) CreateRoleAdmin(ctx context.Context, req entity.RoleEntity) (int64, error) {
	var roleId int64

	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleIdCreated, err := r.repo.CreateRole(txCtx, req)
		if err != nil {
			return err
		}

		if err := r.cacheRole.DeleteRoleCache(ctx, roleIdCreated); err != nil {
			return err
		}

		roleId = roleIdCreated

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] CreateRoleAdmin: %v", err)
		return 0, err
	}

	return roleId, nil
}

// DeleteRoleAdmin implements RoleServiceInterface.
func (r *roleService) DeleteRoleAdmin(ctx context.Context, id int64) error {
	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := r.repo.DeleteRole(txCtx, id); err != nil {
			return err
		}

		if err := r.cacheRole.DeleteRoleCache(ctx, id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] DeleteRoleAdmin: %v", err)
		return err
	}

	return nil
}

// GetRolesAllAdmin implements RoleServiceInterface.
func (r *roleService) GetRolesAllAdmin(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	roles := []entity.RoleEntity{}

	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntities, err := r.repo.GetRolesAll(txCtx, search)
		if err != nil {
			return nil
		}

		roles = roleEntities

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] GetRolesAllAdmin: %v", err)
		return nil, err
	}

	return roles, nil
}

// GetRoleByIdAdmin implements RoleServiceInterface.
func (r *roleService) GetRoleByIdAdmin(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	var (
		role *entity.RoleEntity
	)

	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := r.cacheRole.GetRoleById(txCtx, id)
		if err != nil {
			return err
		}

		role = roleEntity

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] GetRoleByIdAdmin: %v", err)
		return nil, err
	}

	return role, nil
}

// UpdateRoleAdmin implements RoleServiceInterface.
func (r *roleService) UpdateRoleAdmin(ctx context.Context, req entity.RoleEntity) error {
	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := r.repo.UpdateRole(txCtx, req); err != nil {
			return err
		}

		if err := r.cacheRole.DeleteRoleCache(ctx, req.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] UpdateRoleAdmin: %v", err)
		return err
	}

	return nil
}
