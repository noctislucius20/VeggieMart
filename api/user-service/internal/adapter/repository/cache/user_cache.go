package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"
	"user-service/internal/adapter/repository"
	"user-service/internal/core/domain/entity"
	"user-service/utils"
	"user-service/utils/helper"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type UserCacheInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	GetCustomerById(ctx context.Context, id int64) (*entity.UserEntity, error)
	GetProfileById(ctx context.Context, id int64) (*entity.UserEntity, error)
	GetDataByToken(ctx context.Context, token string) (*entity.VerificationTokenEntity, error)
	SetUserSession(ctx context.Context, session entity.SessionEntity) error
	DeleteUserCache(ctx context.Context, id int64) error
}

type userCache struct {
	redisClient *redis.Client
	repoUser    repository.UserRepositoryInterface
	repoToken   repository.VerificationTokenRepositoryInterface
	logger      *log.Logger
}

func NewUserCache(redisClient *redis.Client, repoUser repository.UserRepositoryInterface, repoToken repository.VerificationTokenRepositoryInterface, logger *log.Logger) UserCacheInterface {
	return &userCache{
		redisClient: redisClient,
		repoUser:    repoUser,
		repoToken:   repoToken,
		logger:      logger,
	}
}

// DeleteUserCache implements [UserCacheInterface].
func (u *userCache) DeleteUserCache(ctx context.Context, id int64) error {
	var (
		verifyToken entity.VerificationTokenEntity
		session     entity.SessionEntity
		user        entity.UserEntity
		keys        = []string{
			fmt.Sprintf("user:id:%d:verifytoken", id),
			fmt.Sprintf("user:id:%d:session", id),
			fmt.Sprintf("user:id:%d:email", id),
		}
		delKeys = []string{
			fmt.Sprintf("user:customer:%d", id),
			fmt.Sprintf("user:profile:%d", id),
		}
	)

	if err := helper.RedisBulkGet(ctx, u.redisClient, keys, &verifyToken, &session, &user); err != nil {
		u.logger.Errorf("[UserCache] DeleteUserCache: %v", err)
		return err
	}

	delKeys = helper.AppendKeyIfNotEmpty(delKeys, "user:email:%s", user.Email)
	delKeys = helper.AppendKeyIfNotEmpty(delKeys, "user:verifytoken:%s", verifyToken.Token)
	delKeys = helper.AppendKeyIfNotEmpty(delKeys, "user:session:%s", session.Token)
	delKeys = append(delKeys, keys...)

	if err := u.redisClient.Del(ctx, delKeys...).Err(); err != nil {
		u.logger.Errorf("[UserCache] DeleteUserCache: %v", err)
		return err
	}

	return nil
}

// GetProfileById implements [UserCacheInterface].
func (u *userCache) GetProfileById(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var (
		profile entity.UserEntity
		key     = fmt.Sprintf("user:profile:%d", id)
	)

	// Check redis if data exists.
	val, err := u.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := utils.ErrDataNotFound
			u.logger.Errorf("[UserCache] GetProfileById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &profile)

		return &profile, nil
	}

	profileEntity, err := u.repoUser.GetProfileById(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrDataNotFound) {
			if err := u.redisClient.Set(ctx, key, "null", 10*time.Minute); err != nil {
				u.logger.Errorf("[UserCache] GetProfileById: %v", err)
			}
		}

		u.logger.Errorf("[UserCache] GetProfileById: %v", err)
		return nil, err
	}

	profile = *profileEntity

	// Save to redis
	jsonData, _ := json.Marshal(profile)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := u.redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] GetProfileById: %v", err)
	}

	return &profile, nil
}

// GetCustomerById implements [UserCacheInterface].
func (u *userCache) GetCustomerById(ctx context.Context, id int64) (*entity.UserEntity, error) {
	var (
		user entity.UserEntity
		key  = fmt.Sprintf("user:customer:%d", id)
	)

	// Check redis if data exists.
	val, err := u.redisClient.Get(ctx, key).Result()
	if err == nil {
		// if key exists but value null, return data not found error
		if val == "null" {
			err := utils.ErrDataNotFound
			u.logger.Errorf("[UserCache] GetCustomerById: %v", err)
			return nil, err
		}

		json.Unmarshal([]byte(val), &user)

		return &user, nil
	}

	userEntity, err := u.repoUser.GetCustomerById(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrDataNotFound) {
			if err := u.redisClient.Set(ctx, key, "null", 10*time.Minute); err != nil {
				u.logger.Errorf("[UserCache] GetCustomerById: %v", err)
			}
		}

		u.logger.Errorf("[UserCache] GetCustomerById: %v", err)
		return nil, err
	}

	user = *userEntity

	// Save to redis
	jsonData, _ := json.Marshal(user)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	if err := u.redisClient.Set(ctx, key, jsonData, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] GetCustomerById: %v", err)
	}

	return &user, nil
}

// SetUserSession implements [UserCacheInterface].
func (u *userCache) SetUserSession(ctx context.Context, session entity.SessionEntity) error {
	var (
		key = fmt.Sprintf("user:session:%s", session.Token)
	)

	sessionDataJson, err := json.Marshal(session)
	if err != nil {
		u.logger.Errorf("[UserCache] SetUserSession: %v", err)
		return err
	}

	ttl := 23*time.Hour + time.Duration(rand.Intn(120))*time.Second

	pipe := u.redisClient.Pipeline()

	if err := pipe.Set(ctx, key, session.UserID, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] SetUserSession: %v", err)
		return err
	}
	if err := pipe.Set(ctx, fmt.Sprintf("user:id:%d:session", session.UserID), sessionDataJson, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] SetUserSession: %v", err)
		return err
	}
	if _, err = pipe.Exec(ctx); err != nil {
		u.logger.Errorf("[UserCache] SetUserSession: %v", err)
		return err
	}

	return nil
}

// GetDataByToken implements [UserCacheInterface].
func (u *userCache) GetDataByToken(ctx context.Context, token string) (*entity.VerificationTokenEntity, error) {
	var (
		tokenData entity.VerificationTokenEntity
		key       = fmt.Sprintf("user:verifytoken:%s", token)
	)

	// Check redis if data exists.
	if valKey, err := u.redisClient.Get(ctx, key).Result(); err == nil {
		// if key exists but value null, return data not found error
		if valKey == "null" {
			err := utils.ErrDataNotFound
			u.logger.Errorf("[UserCache] GetDataByToken: %v", err)

			return nil, err
		}
		if val, err := u.redisClient.Get(ctx, fmt.Sprintf("user:id:%s:verifytoken", valKey)).Result(); err == nil {
			json.Unmarshal([]byte(val), &tokenData)

			return &tokenData, nil
		}
	}

	tokenEntity, err := u.repoToken.GetDataByToken(ctx, token)
	if err != nil {
		if errors.Is(err, utils.ErrTokenInvalid) {
			if err := u.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
				u.logger.Errorf("[UserCache] GetDataByToken: %v", err)
			}
		}

		u.logger.Errorf("[UserCache] GetDataByToken: %v", err)
		return nil, err
	}

	tokenData = *tokenEntity

	// Save to redis
	jsonData, _ := json.Marshal(tokenData)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	pipe := u.redisClient.Pipeline()
	if err := pipe.Set(ctx, key, tokenData.UserID, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] GetDataByToken: %v", err)
	}
	if err := pipe.Set(ctx, fmt.Sprintf("user:id:%d:verifytoken", tokenData.UserID), jsonData, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] GetDataByToken: %v", err)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		u.logger.Errorf("[UserCache] GetDataByToken: %v", err)
	}

	return &tokenData, nil
}

// GetUserByEmail implements [UserCacheInterface].
func (u *userCache) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	var (
		user entity.UserEntity
		key  = fmt.Sprintf("user:email:%s", email)
	)

	// Check redis if data exists.
	if valKey, err := u.redisClient.Get(ctx, key).Result(); err == nil {
		// if key exists but value null, return data not found error
		if valKey == "null" {
			err := utils.ErrDataNotFound
			u.logger.Errorf("[UserCache] GetUserByEmail: %v", err)

			return nil, err
		}
		if val, err := u.redisClient.Get(ctx, fmt.Sprintf("user:id:%s:email", valKey)).Result(); err == nil {
			json.Unmarshal([]byte(val), &user)

			return &user, nil
		}
	}

	userEntity, err := u.repoUser.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, utils.ErrDataNotFound) {
			if err := u.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
				u.logger.Errorf("[UserCache] GetUserByEmail: %v", err)
			}
		}

		u.logger.Errorf("[UserCache] GetUserByEmail: %v", err)
		return nil, err
	}

	user = *userEntity

	// Save to redis
	jsonData, _ := json.Marshal(user)
	ttl := 10*time.Minute + time.Duration(rand.Intn(120))*time.Second
	pipe := u.redisClient.Pipeline()
	if err := pipe.Set(ctx, key, user.ID, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] GetUserByEmail: %v", err)
	}
	if err := pipe.Set(ctx, fmt.Sprintf("user:id:%d:email", user.ID), jsonData, ttl).Err(); err != nil {
		u.logger.Errorf("[UserCache] GetUserByEmail: %v", err)
	}
	if _, err = pipe.Exec(ctx); err != nil {
		u.logger.Errorf("[UserCache] GetUserByEmail: %v", err)
	}

	return &user, nil
}
