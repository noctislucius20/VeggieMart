package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils"
	"user-service/utils/conv"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type UserCacheInterface interface {
	SignInByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	SignInSuccess(ctx context.Context, token string, userEntity *entity.UserEntity) error
	VerifyUserSuccess(ctx context.Context, token string, tokenEntity *entity.VerificationTokenEntity) error
}

type userCache struct {
	redisClient *redis.Client
	repoUser    repository.UserRepositoryInterface
	logger      *log.Logger
}

func NewUserCache(redisClient *redis.Client, repoUser repository.UserRepositoryInterface, logger *log.Logger) UserCacheInterface {
	return &userCache{
		redisClient: redisClient,
		repoUser:    repoUser,
		logger:      logger,
	}
}

// VerifyUserSuccess implements [UserCacheInterface].
func (u *userCache) VerifyUserSuccess(ctx context.Context, token string, tokenEntity *entity.VerificationTokenEntity) error {
	var (
		key         = fmt.Sprintf("signin:token:%s", token)
		sessionData = map[string]any{
			"user_id":    tokenEntity.UserID,
			"name":       tokenEntity.User.Name,
			"email":      tokenEntity.User.Email,
			"logged_in":  true,
			"created_at": time.Now().String(),
			"token":      token,
			"role_name":  tokenEntity.User.RoleName,
		}
	)

	sessionDataJson, err := conv.ToJSON(sessionData)
	if err != nil {
		u.logger.Errorf("[UserCache-1] VerifyUserSuccess: %v", err)
		return err
	}

	if err := u.redisClient.Set(ctx, key, sessionDataJson, time.Hour*23).Err(); err != nil {
		u.logger.Errorf("[UserCache-2] VerifyUserSuccess: %v", err)
		return err
	}

	return nil
}

// SignInSuccess implements [UserCacheInterface].
func (u *userCache) SignInSuccess(ctx context.Context, token string, userEntity *entity.UserEntity) error {
	var (
		key         = fmt.Sprintf("signin:token:%s", token)
		sessionData = map[string]any{
			"user_id":    userEntity.ID,
			"name":       userEntity.Name,
			"email":      userEntity.Email,
			"logged_in":  true,
			"created_at": time.Now().String(),
			"token":      token,
			"role_name":  userEntity.RoleName,
		}
	)

	sessionDataJson, err := conv.ToJSON(sessionData)
	if err != nil {
		u.logger.Errorf("[UserCache-1] SignInSuccess: %v", err)
		return err
	}

	if err := u.redisClient.Set(ctx, key, sessionDataJson, time.Hour*23).Err(); err != nil {
		u.logger.Errorf("[UserCache-2] SignInSuccess: %v", err)
		return err
	}

	return nil
}

// SignInByEmail implements [UserCacheInterface].
func (u *userCache) SignInByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	var (
		user entity.UserEntity
		key  = fmt.Sprintf("signin:email:%s", email)
	)

	// Check redis if data exists.
	val, err := u.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := errors.New(utils.DATA_NOT_FOUND)
			u.logger.Errorf("[UserCache-1] SignInByEmail: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &user)

		return &user, nil
	}

	userEntity, err := u.repoUser.GetUserByEmail(ctx, email)
	if err != nil {
		if err.Error() == utils.DATA_NOT_FOUND {
			if err := u.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
				u.logger.Errorf("[UserCache-1] SignInByEmail: %v", err)
			}
		}

		u.logger.Errorf("[UserCache-2] SignInByEmail: %v", err)
		return nil, err
	}

	user = *userEntity

	// Save to redis
	jsonData, _ := json.Marshal(user)
	if err := u.redisClient.Set(ctx, key, jsonData, 1*time.Minute).Err(); err != nil {
		u.logger.Errorf("[UserCache-3] SignInByEmail: %v", err)
	}

	return &user, nil
}
