package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"user-service/config"
	"user-service/internal/adapter/repository"
	"user-service/internal/adapter/repository/cache"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service/transaction"
	"user-service/utils"
	"user-service/utils/conv"

	"github.com/google/uuid"
	"github.com/labstack/gommon/log"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) (int64, error)
	ForgotPassword(ctx context.Context, req entity.UserEntity) error
	VerifyToken(ctx context.Context, token string) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, req entity.UserEntity) error
	GetProfileById(ctx context.Context, userId int64) (*entity.UserEntity, error)
	UpdateProfile(ctx context.Context, req entity.UserEntity) error

	// Admin customer management functions can be added here
	GetBatchCustomers(ctx context.Context, userIds []int64) ([]entity.UserEntity, error)
	GetCustomersAll(ctx context.Context, query entity.QueryStringEntity) ([]entity.UserEntity, int64, int64, error)
	GetCustomerById(ctx context.Context, customerId int64) (*entity.UserEntity, error)
	CreateCustomer(ctx context.Context, req entity.UserEntity) (int64, error)
	UpdateCustomer(ctx context.Context, req entity.UserEntity) error
	DeleteCustomer(ctx context.Context, customerId int64) error
}

type userService struct {
	repo        repository.UserRepositoryInterface
	repoOutbox  repository.OutboxEventInterface
	cfg         *config.Config
	jwtService  JwtServiceInterface
	repoToken   repository.VerificationTokenRepositoryInterface
	cacheUser   cache.UserCacheInterface
	roleService RoleServiceInterface
	txManager   transaction.TransactionManager
	logger      *log.Logger
}

func NewUserService(repo repository.UserRepositoryInterface, cfg *config.Config, jwtService JwtServiceInterface, repoToken repository.VerificationTokenRepositoryInterface, repoOutbox repository.OutboxEventInterface, roleService RoleServiceInterface, cacheUser cache.UserCacheInterface, txManager transaction.TransactionManager, logger *log.Logger) UserServiceInterface {
	return &userService{
		repo:        repo,
		cfg:         cfg,
		jwtService:  jwtService,
		repoOutbox:  repoOutbox,
		repoToken:   repoToken,
		cacheUser:   cacheUser,
		roleService: roleService,
		txManager:   txManager,
		logger:      logger,
	}
}

// GetBatchCustomers implements [UserServiceInterface].
func (u *userService) GetBatchCustomers(ctx context.Context, userIds []int64) ([]entity.UserEntity, error) {
	users := []entity.UserEntity{}

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		usersEntity, err := u.repo.GetBatchCustomers(txCtx, userIds)
		if err != nil {
			return err
		}

		users = usersEntity

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] GetBatchCustomers: %v", err)
		return nil, err
	}

	return users, nil
}

// DeleteCustomer implements UserServiceInterface.
func (u *userService) DeleteCustomer(ctx context.Context, customerId int64) error {
	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := u.repo.DeleteCustomer(txCtx, customerId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] DeleteCustomer: %v", err)
		return err
	}

	// redis delete key
	// key := fmt.Sprintf("customer:%d", customerId)
	// if err := u.redisClient.Del(ctx, key).Err(); err != nil {
	// 	u.logger.Errorf("[UserService-2] DeleteCustomer: %v", err)
	// }

	return nil
}

// UpdateCustomer implements UserServiceInterface.
func (u *userService) UpdateCustomer(ctx context.Context, req entity.UserEntity) error {
	var password string

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if req.Password != "" {
			password = req.Password
			hashedPassword, err := conv.HashPassword(req.Password)
			if err != nil {
				return err
			}
			req.Password = hashedPassword
		}

		roleEntity, err := u.roleService.GetRoleByIdAdmin(txCtx, 2)
		if err != nil {
			return err
		}
		req.RoleId = roleEntity.ID

		if err := u.repo.UpdateCustomer(txCtx, req); err != nil {
			return err
		}

		if password != "" {
			payloadMessage := fmt.Sprintf("Your password has been updated in Sayur App. Please use this credential to login: \n Email: %s\nPassword: %s", req.Email, password)

			publishEmailPayload := map[string]any{
				"receiver_email":    req.Email,
				"message":           payloadMessage,
				"subject":           "Updated Data Customer",
				"receiver_id":       req.ID,
				"notification_type": "EMAIL",
			}
			if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_UPDATE_CUSTOMER, publishEmailPayload, &req.ID); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] UpdateCustomer: %v", err)
		return err
	}

	// redis delete key
	// key := fmt.Sprintf("customer:%d", req.ID)
	// if err := u.redisClient.Del(ctx, key).Err(); err != nil {
	// 	u.logger.Errorf("[UserService-2] UpdateCustomer: %v", err)
	// }

	return nil
}

// CreateCustomer implements UserServiceInterface.
func (u *userService) CreateCustomer(ctx context.Context, req entity.UserEntity) (int64, error) {
	var customerId int64

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		password := req.Password

		hashedPassword, err := conv.HashPassword(password)
		if err != nil {
			return err
		}

		req.Password = hashedPassword

		roleEntity, err := u.roleService.GetRoleByIdAdmin(txCtx, 2)
		if err != nil {
			return err
		}
		req.RoleId = roleEntity.ID

		customer, err := u.repo.CreateCustomer(txCtx, req)
		if err != nil {
			return err
		}

		payloadMessage := fmt.Sprintf(`
			<p>You have been registered in Sayur App. Please use this credential to login:</p>
			<p><b>Email: %s </b></p>
			<p><b>Password: %s </b></p>`, req.Email, password)

		publishEmailPayload := map[string]any{
			"receiver_email":    req.Email,
			"message":           payloadMessage,
			"subject":           "Verify Your Account",
			"receiver_id":       customerId,
			"notification_type": "EMAIL",
		}
		if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_CREATE_CUSTOMER, publishEmailPayload, &customerId); err != nil {
			return err
		}

		customerId = customer

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] CreateCustomer: %v", err)
		return 0, err
	}

	// redis delete key
	// key := fmt.Sprintf("customer:%d", req.ID)
	// if err := u.redisClient.Del(ctx, key).Err(); err != nil {
	// 	u.logger.Errorf("[UserService-2] CreateCustomer: %v", err)
	// }

	return customerId, nil
}

// GetCustomerById implements UserServiceInterface.
func (u *userService) GetCustomerById(ctx context.Context, customerId int64) (*entity.UserEntity, error) {
	var (
		customer entity.UserEntity
		// key      = fmt.Sprintf("customer:%d", customerId)
	)

	// Check redis if data exists.
	// val, err := u.redisClient.Get(ctx, key).Result()
	// if err == nil {
	// 	// if key exists but value null, return data not found error
	// 	if val == "null" {
	// 		err := errors.New(utils.DATA_NOT_FOUND)
	// 		return nil, err
	// 	}

	// 	json.Unmarshal([]byte(val), &customer)
	// 	return &customer, nil
	// }

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		customerEntity, err := u.repo.GetCustomerById(txCtx, customerId)
		if err != nil {
			return err
		}

		customer = *customerEntity

		return nil
	}); err != nil {
		// Save to redis (create key with null value if data not found)
		// if err.Error() == utils.DATA_NOT_FOUND {
		// 	if err := u.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
		// 		u.logger.Errorf("[UserService-1] GetCustomerById: %v", err)
		// 	}
		// }

		u.logger.Errorf("[UserService-2] GetCustomerById: %v", err)
		return nil, err
	}

	// Save to redis
	// jsonData, _ := json.Marshal(customer)
	// if err := u.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
	// 	u.logger.Errorf("[UserService-3] GetCustomerById: %v", err)
	// }

	return &customer, nil
}

// GetCustomersAll implements UserServiceInterface.
func (u *userService) GetCustomersAll(ctx context.Context, query entity.QueryStringEntity) ([]entity.UserEntity, int64, int64, error) {
	var (
		customers  []entity.UserEntity
		countData  int64
		totalPages int64
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		customerEntities, count, pages, err := u.repo.GetAllCustomers(txCtx, query)
		if err != nil {
			return err
		}

		customers, countData, totalPages = customerEntities, count, pages

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] GetCustomersAll: %v", err)
		return nil, 0, 0, err
	}

	return customers, countData, totalPages, nil
}

// UpdateProfile implements UserServiceInterface.
func (u *userService) UpdateProfile(ctx context.Context, req entity.UserEntity) error {
	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := u.repo.UpdateProfile(txCtx, req); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] UpdateProfile: %v", err)
		return err
	}

	// redis delete key
	// key := fmt.Sprintf("user:%d", req.ID)
	// if err := u.redisClient.Del(ctx, key).Err(); err != nil {
	// 	u.logger.Errorf("[UserService-2] UpdateProfile: %v", err)
	// }

	return nil
}

// GetProfileById implements UserServiceInterface.
func (u *userService) GetProfileById(ctx context.Context, userId int64) (*entity.UserEntity, error) {
	var (
		profile entity.UserEntity
		// key     = fmt.Sprintf("user:%d", userId)
	)

	// Check redis if data exists.
	// val, err := u.redisClient.Get(ctx, key).Result()
	// if err == nil {
	// 	// if key exists but value null, return data not found error
	// 	if val == "null" {
	// 		err := errors.New(utils.DATA_NOT_FOUND)
	// 		return nil, err
	// 	}

	// 	json.Unmarshal([]byte(val), &profile)
	// 	return &profile, nil
	// }

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		profileEntity, err := u.repo.GetProfileById(txCtx, userId)
		if err != nil {
			return err
		}

		profile = *profileEntity

		return nil
	}); err != nil {
		// Save to redis (create key with null value if data not found)
		// if err.Error() == utils.DATA_NOT_FOUND {
		// 	if err := u.redisClient.Set(ctx, key, "null", 1*time.Minute); err != nil {
		// 		u.logger.Errorf("[UserService-1] GetProfileById: %v", err)
		// 	}
		// }

		u.logger.Errorf("[UserService-2] GetProfileById: %v", err)
		return nil, err
	}

	// Save to redis
	// jsonData, _ := json.Marshal(profile)
	// if err := u.redisClient.Set(ctx, key, jsonData, 1*time.Hour).Err(); err != nil {
	// 	u.logger.Errorf("[UserService-3] GetProfileById: %v", err)
	// }

	return &profile, nil
}

// UpdatePassword implements UserServiceInterface.
func (u *userService) UpdatePassword(ctx context.Context, req entity.UserEntity) error {
	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tokenEntity, err := u.repoToken.GetDataByToken(txCtx, req.Token)
		if err != nil {
			return err
		}

		if time.Now().After(tokenEntity.ExpiresAt) {
			err := errors.New(utils.TOKEN_EXPIRED)
			if err := u.repoToken.DeleteVerificationToken(txCtx, tokenEntity.ID); err != nil {
				return err
			}
			return err
		}

		if tokenEntity.TokenType != utils.NOTIF_EMAIL_FORGOT_PASSWORD {
			err := errors.New(utils.TOKEN_INVALID)
			return err
		}

		hashedPassword, err := conv.HashPassword(req.Password)
		if err != nil {
			return err
		}

		req.Password = hashedPassword
		req.ID = tokenEntity.UserID

		if err := u.repo.UpdatePasswordById(txCtx, req); err != nil {
			return err
		}

		if err := u.repoToken.DeleteVerificationToken(txCtx, tokenEntity.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] UpdatePassword: %v", err)
		return err
	}

	return nil
}

// VerifyToken implements UserServiceInterface.
func (u *userService) VerifyToken(ctx context.Context, token string) (*entity.UserEntity, error) {
	var user *entity.UserEntity

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tokenEntity, err := u.repoToken.GetDataByToken(txCtx, token)
		if err != nil {
			return err
		}

		if tokenEntity.TokenType != utils.NOTIF_EMAIL_VERIFICATION {
			err := errors.New(utils.TOKEN_INVALID)
			return err
		}

		if time.Now().After(tokenEntity.ExpiresAt) {
			err := errors.New(utils.TOKEN_EXPIRED)
			if err := u.repoToken.DeleteVerificationToken(txCtx, tokenEntity.ID); err != nil {
				return err
			}
			return err
		}

		if err := u.repo.UpdateUserVerified(txCtx, tokenEntity.UserID); err != nil {
			return err
		}

		if err := u.repoToken.DeleteVerificationToken(txCtx, tokenEntity.ID); err != nil {
			return err
		}

		accessToken, err := u.jwtService.GenerateToken(tokenEntity.UserID)
		if err != nil {
			return err
		}

		// sessionData := map[string]any{
		// 	"user_id":    tokenEntity.UserID,
		// 	"name":       tokenEntity.User.Name,
		// 	"email":      tokenEntity.User.Email,
		// 	"logged_in":  true,
		// 	"created_at": time.Now().String(),
		// 	"token":      accessToken,
		// 	"role_name":  tokenEntity.User.RoleName,
		// }
		// sessionDataJson, _ := conv.ToJSON(sessionData)

		// if err := u.redisClient.Set(ctx, accessToken, sessionDataJson, time.Hour*23).Err(); err != nil {
		// 	return err
		// }

		tokenEntity.User.Token = accessToken

		user = &tokenEntity.User

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] VerifyToken: %v", err)
		return nil, err
	}

	return user, nil

}

// ForgotPassword implements UserServiceInterface.
func (u *userService) ForgotPassword(ctx context.Context, req entity.UserEntity) error {
	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		user, err := u.repo.GetUserByEmail(txCtx, req.Email)
		if err != nil {
			return err
		}

		if user.IsVerified == false {
			err := errors.New(utils.EMAIL_NOT_VERIFIED)
			return err
		}

		token := uuid.New().String()
		reqEntity := entity.VerificationTokenEntity{
			UserID:    user.ID,
			Token:     token,
			TokenType: utils.NOTIF_EMAIL_FORGOT_PASSWORD,
		}

		if err := u.repoToken.CreateVerificationToken(txCtx, reqEntity); err != nil {
			return err
		}

		urlForgot := fmt.Sprintf("%s/reset-password?token=%s", u.cfg.App.UrlUsersService, token)
		payloadMessage := fmt.Sprintf("Please click link below to reset your password: %v", urlForgot)

		publishEmailPayload := map[string]any{
			"receiver_email":    req.Email,
			"message":           payloadMessage,
			"subject":           "Reset Password",
			"receiver_id":       user.ID,
			"notification_type": "EMAIL",
		}
		if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_FORGOT_PASSWORD, publishEmailPayload, &user.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] ForgotPassword: %v", err)
		return err
	}

	return nil
}

// CreateUserAccount implements UserServiceInterface.
func (u *userService) CreateUserAccount(ctx context.Context, req entity.UserEntity) (int64, error) {
	var userId int64

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		password, err := conv.HashPassword(req.Password)
		if err != nil {
			return err
		}

		roleEntity, err := u.roleService.GetRoleByIdAdmin(txCtx, 2)
		if err != nil {
			return err
		}

		req.Password = password
		req.Token = uuid.New().String()
		req.RoleId = roleEntity.ID

		userIdCreated, err := u.repo.CreateUserAccount(txCtx, req)
		if err != nil {
			return err
		}

		reqEntity := entity.VerificationTokenEntity{
			UserID:    userIdCreated,
			Token:     req.Token,
			TokenType: utils.NOTIF_EMAIL_VERIFICATION,
		}

		if err := u.repoToken.CreateVerificationToken(txCtx, reqEntity); err != nil {
			return err
		}

		urlVerify := fmt.Sprintf("%s/verify?token=%s", u.cfg.App.UrlUsersService, req.Token)
		payloadMessage := fmt.Sprintf("Please click link below to activate your account: %v", urlVerify)

		publishEmailPayload := map[string]any{
			"receiver_email":    req.Email,
			"message":           payloadMessage,
			"subject":           "Account Exists",
			"receiver_id":       userIdCreated,
			"notification_type": "EMAIL",
		}
		if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_VERIFICATION, publishEmailPayload, &userIdCreated); err != nil {
			return err
		}

		userId = userIdCreated

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] CreateUserAccount: %v", err)
		return 0, err
	}

	return userId, nil
}

// SignIn implements UserServiceInterface.
func (u *userService) SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error) {
	var (
		user  *entity.UserEntity
		token string
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		userEntity, err := u.cacheUser.SignInByEmail(txCtx, req.Email)
		if err != nil {
			return err
		}

		if checkPass := conv.CheckPasswordHash(req.Password, userEntity.Password); !checkPass {
			err = errors.New(utils.LOGIN_INVALID)
			return err
		}

		if userEntity.IsVerified == false {
			err := errors.New(utils.EMAIL_NOT_VERIFIED)
			return err
		}

		tokenString, err := u.jwtService.GenerateToken(userEntity.ID)
		if err != nil {
			return err
		}

		if err := u.cacheUser.SignInSuccess(txCtx, tokenString, userEntity); err != nil {
			return err
		}

		user, token = userEntity, tokenString

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService-1] SignIn: %v", err)
		return nil, "", err
	}

	return user, token, nil
}
