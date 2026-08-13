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

	"github.com/labstack/gommon/log"
)

type UserServiceInterface interface {
	SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error)
	SignUp(ctx context.Context, req entity.UserEntity) (int64, error)
	ForgotPassword(ctx context.Context, req entity.UserEntity) error
	ActivateAccount(ctx context.Context, token string) (*entity.UserEntity, error)
	UpdatePassword(ctx context.Context, req entity.UserEntity) error
	GetProfileById(ctx context.Context, userId int64) (*entity.UserEntity, error)
	UpdateProfile(ctx context.Context, req entity.UserEntity) (string, string, error)

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
		u.logger.Errorf("[UserService] GetBatchCustomers: %v", err)
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

		if err := u.cacheUser.DeleteUserCache(ctx, customerId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] DeleteCustomer: %v", err)
		return err
	}

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

		_, err := u.roleService.GetRoleByIdAdmin(txCtx, req.RoleID)
		if err != nil {
			return utils.ErrInternalServerError
		}

		if err := u.repo.UpdateCustomer(txCtx, req); err != nil {
			return err
		}

		if password != "" {
			payloadMessage := fmt.Sprintf("Your password has been updated in Sayur App. Please use this credential to login: \n Email: %s\nPassword: %s", req.Email, password)

			publishEmailPayload := map[string]any{
				"receiver_email":      req.Email,
				"message":             payloadMessage,
				"subject":             "Updated Data Customer",
				"receiver_id":         req.ID,
				"notification_method": "EMAIL",
			}

			if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_UPDATE_CUSTOMER, publishEmailPayload, &req.ID); err != nil {
				return err
			}
		}

		// delete all user cache
		if err := u.cacheUser.DeleteUserCache(txCtx, req.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] UpdateCustomer: %v", err)
		return err
	}

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

		if _, err := u.roleService.GetRoleByIdAdmin(txCtx, req.RoleID); err != nil {
			return utils.ErrInternalServerError
		}

		customerIdCreated, err := u.repo.CreateCustomer(txCtx, req)
		if err != nil {
			return err
		}

		payloadMessage := fmt.Sprintf(`
			<p>You have been registered in Sayur App. Please use this credential to login:</p>
			<p><b>Email: %s </b></p>
			<p><b>Password: %s </b></p>`, req.Email, password)

		publishEmailPayload := map[string]any{
			"receiver_email":      req.Email,
			"message":             payloadMessage,
			"subject":             "Verify Your Account",
			"receiver_id":         customerIdCreated,
			"notification_method": "EMAIL",
		}

		if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_CREATE_CUSTOMER, publishEmailPayload, &customerId); err != nil {
			return err
		}

		if err := u.cacheUser.DeleteUserCache(ctx, customerIdCreated); err != nil {
			return err
		}

		customerId = customerIdCreated

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] CreateCustomer: %v", err)
		return 0, err
	}

	return customerId, nil
}

// GetCustomerById implements UserServiceInterface.
func (u *userService) GetCustomerById(ctx context.Context, customerId int64) (*entity.UserEntity, error) {
	var (
		customer entity.UserEntity
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		customerEntity, err := u.cacheUser.GetCustomerById(txCtx, customerId)
		if err != nil {
			return err
		}

		customer = *customerEntity

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] GetCustomerById: %v", err)
		return nil, err
	}

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
		u.logger.Errorf("[UserService] GetCustomersAll: %v", err)
		return nil, 0, 0, err
	}

	return customers, countData, totalPages, nil
}

// UpdateProfile implements UserServiceInterface.
func (u *userService) UpdateProfile(ctx context.Context, req entity.UserEntity) (string, string, error) {
	var (
		token    string
		roleName string
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := u.repo.UpdateProfile(txCtx, req); err != nil {
			return err
		}

		profile, err := u.cacheUser.GetProfileById(txCtx, req.ID)
		if err != nil {
			return utils.ErrInternalServerError
		}

		tokenString, err := u.jwtService.GenerateToken(req.ID)
		if err != nil {
			return err
		}

		session := entity.SessionEntity{
			UserID:    req.ID,
			Name:      req.Name,
			Email:     req.Email,
			LoggedIn:  true,
			CreatedAt: time.Now().String(),
			Token:     tokenString,
			RoleID:    profile.RoleID,
		}

		if err := u.cacheUser.DeleteUserCache(txCtx, req.ID); err != nil {
			return err
		}

		if err := u.cacheUser.SetUserSession(txCtx, session); err != nil {
			return err
		}

		token = tokenString
		roleName = profile.RoleName

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] UpdateProfile: %v", err)
		return "", "", err
	}

	return token, roleName, nil
}

// GetProfileById implements UserServiceInterface.
func (u *userService) GetProfileById(ctx context.Context, userId int64) (*entity.UserEntity, error) {
	var (
		profile entity.UserEntity
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		profileEntity, err := u.cacheUser.GetProfileById(txCtx, userId)
		if err != nil {
			return err
		}

		roleEntity, err := u.roleService.GetRoleByIdAdmin(txCtx, profileEntity.RoleID)
		if err != nil {
			return utils.ErrInternalServerError
		}

		profileEntity.RoleName = roleEntity.Name

		profile = *profileEntity

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] GetProfileById: %v", err)
		return nil, err
	}

	return &profile, nil
}

// UpdatePassword implements UserServiceInterface.
func (u *userService) UpdatePassword(ctx context.Context, req entity.UserEntity) error {
	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tokenEntity, err := u.cacheUser.GetDataByToken(txCtx, req.Token)
		if err != nil {
			return err
		}

		if tokenEntity.TokenType != utils.NOTIF_EMAIL_FORGOT_PASSWORD {
			err := utils.ErrTokenInvalid
			return err
		}

		if time.Now().After(tokenEntity.ExpiresAt) {
			err := utils.ErrTokenExpired
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

		if err := u.cacheUser.DeleteUserCache(ctx, req.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] UpdatePassword: %v", err)
		return err
	}

	return nil
}

// ActivateAccount implements UserServiceInterface.
func (u *userService) ActivateAccount(ctx context.Context, token string) (*entity.UserEntity, error) {
	var (
		user    *entity.UserEntity
		session entity.SessionEntity
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		tokenEntity, err := u.cacheUser.GetDataByToken(txCtx, token)
		if err != nil {
			return err
		}

		if tokenEntity.TokenType != utils.NOTIF_EMAIL_VERIFICATION {
			err := utils.ErrTokenInvalid
			return err
		}

		if time.Now().After(tokenEntity.ExpiresAt) {
			err := utils.ErrTokenExpired
			return err
		}

		if err := u.repo.UpdateUserVerified(txCtx, tokenEntity.UserID); err != nil {
			return err
		}

		if err := u.repoToken.DeleteVerificationToken(txCtx, tokenEntity.ID); err != nil {
			return err
		}

		if err := u.cacheUser.DeleteUserCache(ctx, tokenEntity.UserID); err != nil {
			return err
		}

		accessToken, err := u.jwtService.GenerateToken(tokenEntity.UserID)
		if err != nil {
			return err
		}

		session = entity.SessionEntity{
			UserID:    tokenEntity.User.ID,
			Name:      tokenEntity.User.Name,
			Email:     tokenEntity.User.Email,
			LoggedIn:  true,
			CreatedAt: time.Now().String(),
			Token:     accessToken,
			RoleID:    tokenEntity.User.RoleID,
		}

		if err := u.cacheUser.SetUserSession(txCtx, session); err != nil {
			return err
		}

		roleEntity, err := u.roleService.GetRoleByIdAdmin(txCtx, tokenEntity.User.RoleID)
		if err != nil {
			return utils.ErrInternalServerError
		}

		tokenEntity.User.Token = accessToken
		tokenEntity.User.RoleName = roleEntity.Name

		user = &tokenEntity.User

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] ActivateAccount: %v", err)
		return nil, err
	}

	return user, nil
}

// ForgotPassword implements UserServiceInterface.
func (u *userService) ForgotPassword(ctx context.Context, req entity.UserEntity) error {
	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		user, err := u.cacheUser.GetUserByEmail(txCtx, req.Email)
		if err != nil {
			return err
		}

		if user.IsVerified == false {
			err := utils.ErrEmailNotVerified
			return err
		}

		reqEntity := entity.VerificationTokenEntity{
			UserID:    user.ID,
			Token:     req.Token,
			TokenType: utils.NOTIF_EMAIL_FORGOT_PASSWORD,
		}

		if err := u.repoToken.CreateVerificationToken(txCtx, reqEntity); err != nil {
			return err
		}

		urlForgot := fmt.Sprintf("%s/auth/reset-password?token=%s", u.cfg.App.UrlFrontEnd, req.Token)
		payloadMessage := fmt.Sprintf("Please click link below to reset your password: %v", urlForgot)

		publishEmailPayload := map[string]any{
			"receiver_email":      req.Email,
			"message":             payloadMessage,
			"subject":             "Reset Password",
			"receiver_id":         user.ID,
			"notification_method": "EMAIL",
		}
		if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_FORGOT_PASSWORD, publishEmailPayload, &user.ID); err != nil {
			return err
		}

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] ForgotPassword: %v", err)
		return err
	}

	return nil
}

// SignUp implements UserServiceInterface.
func (u *userService) SignUp(ctx context.Context, req entity.UserEntity) (int64, error) {
	var userId int64

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		password, err := conv.HashPassword(req.Password)
		if err != nil {
			return err
		}

		if _, err := u.roleService.GetRoleByIdAdmin(txCtx, req.RoleID); err != nil {
			return utils.ErrInternalServerError
		}

		req.Password = password

		userIdCreated, err := u.repo.SignUp(txCtx, req)
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

		urlVerify := fmt.Sprintf("%s/auth/verify-account?token=%s", u.cfg.App.UrlFrontEnd, req.Token)

		payloadMessage := fmt.Sprintf("Please click link below to activate your account: %v", urlVerify)

		publishEmailPayload := map[string]any{
			"receiver_email":      req.Email,
			"message":             payloadMessage,
			"subject":             "Account Exists",
			"receiver_id":         userIdCreated,
			"notification_method": "EMAIL",
		}

		if err := u.repoOutbox.CreateEvent(txCtx, utils.NOTIF_EMAIL_VERIFICATION, publishEmailPayload, &userIdCreated); err != nil {
			return err
		}

		if err := u.cacheUser.DeleteUserCache(ctx, userIdCreated); err != nil {
			return err
		}

		userId = userIdCreated

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] SignUp: %v", err)
		return 0, err
	}

	return userId, nil
}

// SignIn implements UserServiceInterface.
func (u *userService) SignIn(ctx context.Context, req entity.UserEntity) (*entity.UserEntity, string, error) {
	var (
		user    *entity.UserEntity
		session entity.SessionEntity
		token   string
	)

	if err := u.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		userEntity, err := u.cacheUser.GetUserByEmail(txCtx, req.Email)
		if err != nil {
			if errors.Is(err, utils.ErrDataNotFound) {
				err = utils.ErrLoginInvalid
			}
			return err
		}

		if checkPass := conv.CheckPasswordHash(req.Password, userEntity.Password); !checkPass {
			err = utils.ErrLoginInvalid
			return err
		}

		if userEntity.IsVerified == false {
			err := utils.ErrEmailNotVerified
			return err
		}

		tokenString, err := u.jwtService.GenerateToken(userEntity.ID)
		if err != nil {
			return err
		}

		session = entity.SessionEntity{
			UserID:    userEntity.ID,
			Name:      userEntity.Name,
			Email:     userEntity.Email,
			LoggedIn:  true,
			CreatedAt: time.Now().String(),
			Token:     tokenString,
			RoleID:    userEntity.RoleID,
		}

		if err := u.cacheUser.SetUserSession(txCtx, session); err != nil {
			return err
		}

		roleEntity, err := u.roleService.GetRoleByIdAdmin(txCtx, userEntity.RoleID)
		if err != nil {
			return utils.ErrInternalServerError
		}

		userEntity.RoleName = roleEntity.Name

		user, token = userEntity, tokenString

		return nil
	}); err != nil {
		u.logger.Errorf("[UserService] SignIn: %v", err)
		return nil, "", err
	}

	return user, token, nil
}
