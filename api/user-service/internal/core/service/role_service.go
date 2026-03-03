package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service/transaction"
	"user-service/utils"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type RoleServiceInterface interface {
	GetRolesAllAdmin(ctx context.Context, search string) ([]entity.RoleEntity, error)
	GetRoleByIdAdmin(ctx context.Context, id int64) (*entity.RoleEntity, error)
	CreateRoleAdmin(ctx context.Context, req entity.RoleEntity) (int64, error)
	DeleteRoleAdmin(ctx context.Context, id int64) error
	UpdateRoleAdmin(ctx context.Context, req entity.RoleEntity) error
}

type roleService struct {
	repo        repository.RoleRepositoryInterface
	redisClient *redis.Client
	txManager   transaction.TransactionManager
	logger      *log.Logger
}

func NewRoleService(repo repository.RoleRepositoryInterface, redisClient *redis.Client, txManager transaction.TransactionManager, logger *log.Logger) RoleServiceInterface {
	return &roleService{
		repo:        repo,
		redisClient: redisClient,
		txManager:   txManager,
		logger:      logger}
}

// CreateRoleAdmin implements RoleServiceInterface.
func (r *roleService) CreateRoleAdmin(ctx context.Context, req entity.RoleEntity) (int64, error) {
	var roleId int64

	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		role, err := r.repo.CreateRole(txCtx, req)
		if err != nil {
			return err
		}

		roleId = role

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] CreateRoleAdmin: %v", err)
		return 0, err
	}

	// redis delete key
	key := fmt.Sprintf("role:%d", roleId)
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		r.logger.Errorf("[RoleService-2] CreateRoleAdmin: %v", err)
	}

	return roleId, nil
}

// DeleteRoleAdmin implements RoleServiceInterface.
func (r *roleService) DeleteRoleAdmin(ctx context.Context, id int64) error {
	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := r.repo.DeleteRole(txCtx, id); err != nil {
			return err
		}

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] DeleteRoleAdmin: %v", err)
		return err
	}

	// redis delete key
	key := fmt.Sprintf("role:%d", id)
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		r.logger.Errorf("[RoleService-2] DeleteRoleAdmin: %v", err)
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
		key  = fmt.Sprintf("role:%d", id)
	)

	// Check redis if data exists.
	val, err := r.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			r.logger.Errorf("[RoleService-1] GetRoleByIdAdmin: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &role)
		return role, nil
	}

	// Query DB
	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		roleEntity, err := r.repo.GetRoleByIdOrName(txCtx, id, "")
		if err != nil {
			return err
		}

		role = roleEntity

		return nil
	}); err != nil {
		// Save to redis (create key with null value if data not found)
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := r.redisClient.Set(ctx, key, "null", 1*time.Minute).Err(); err != nil {
				r.logger.Errorf("[RoleService-2] GetRoleByIdAdmin: %v", err)
			}
		}

		r.logger.Errorf("[RoleService-3] GetRoleByIdAdmin: %v", err)
		return nil, err
	}

	// Save to redis
	jsonData, _ := json.Marshal(role)
	if err := r.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
		r.logger.Errorf("[RoleService-4] GetRoleByIdAdmin: %v", err)
	}

	return role, nil
}

// UpdateRoleAdmin implements RoleServiceInterface.
func (r *roleService) UpdateRoleAdmin(ctx context.Context, req entity.RoleEntity) error {
	if err := r.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := r.repo.UpdateRole(txCtx, req); err != nil {
			return err
		}

		return nil
	}); err != nil {
		r.logger.Errorf("[RoleService-1] UpdateRoleAdmin: %v", err)
		return err
	}

	// redis delete key
	key := fmt.Sprintf("role:%d", req.ID)
	if err := r.redisClient.Del(ctx, key).Err(); err != nil {
		r.logger.Errorf("[RoleService-2] UpdateRoleAdmin: %v", err)
	}

	return nil
}
