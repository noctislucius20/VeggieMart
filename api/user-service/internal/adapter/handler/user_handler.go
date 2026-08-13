package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"
	middlewareGateway "user-service/internal/middleware"
	"user-service/utils"
	"user-service/utils/conv"
	"user-service/utils/logger"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UserHandlerInterface interface {
	SignIn(c echo.Context) error
	SignUp(c echo.Context) error
	ForgotPassword(c echo.Context) error
	ActivateAccount(c echo.Context) error
	UpdatePassword(c echo.Context) error
	GetProfileById(c echo.Context) error
	UpdateProfile(c echo.Context) error

	// Admin customer management functions can be added here
	GetBatchCustomersAdmin(c echo.Context) error
	GetCustomersAllAdmin(c echo.Context) error
	GetCustomerByIdAdmin(c echo.Context) error
	CreateCustomerAdmin(c echo.Context) error
	UpdateCustomerAdmin(c echo.Context) error
	DeleteCustomerAdmin(c echo.Context) error
}

type userHandler struct {
	userService service.UserServiceInterface
}

func NewUserHandler(e *echo.Echo, userService service.UserServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client) UserHandlerInterface {
	userHandler := &userHandler{userService: userService}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	userGroup := e.Group("/users")
	internalUserGroup := e.Group("/internal/users")

	userGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))
	internalUserGroup.Use(middlewareGateway.InternalServiceMiddleware(cfg))

	userGroup.POST("/signin", userHandler.SignIn)
	userGroup.POST("/signup", userHandler.SignUp)
	userGroup.POST("/forgot-password", userHandler.ForgotPassword)
	userGroup.GET("/activate-account", userHandler.ActivateAccount)
	userGroup.PUT("/reset-password", userHandler.UpdatePassword)

	adminPermission := []string{
		"users:read:all",
		"users:write:all",
		"users:update:all",
		"users:delete:all",
	}

	authPermission := []string{
		"users:read:own",
		"users:write:own",
		"users:update:own",
		"users:delete:own",
	}

	// adminGroup := e.Group("/admin", mid.CheckToken())
	adminCustomerGroup := userGroup.Group("/customers", mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	adminCustomerGroup.GET("", userHandler.GetCustomersAllAdmin)
	adminCustomerGroup.POST("", userHandler.CreateCustomerAdmin)
	adminCustomerGroup.POST("/batch", userHandler.GetBatchCustomersAdmin)
	adminCustomerGroup.GET("/:id", userHandler.GetCustomerByIdAdmin)
	adminCustomerGroup.PUT("/:id", userHandler.UpdateCustomerAdmin)
	adminCustomerGroup.DELETE("/:id", userHandler.DeleteCustomerAdmin)

	// authGroup := e.Group("/auth", mid.CheckToken())
	authCustomerGroup := userGroup.Group("/profile", mid.CheckToken(), mid.RequiredPermission(authPermission...))
	authCustomerGroup.GET("", userHandler.GetProfileById)
	authCustomerGroup.PUT("", userHandler.UpdateProfile)

	internalUserGroup.GET("/profile", userHandler.GetProfileById, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	return userHandler
}

// GetBatchCustomersAdmin implements [UserHandlerInterface].
func (u *userHandler) GetBatchCustomersAdmin(c echo.Context) error {
	var (
		ctx       = c.Request().Context()
		respBatch = []response.CustomerBatchResponse{}
		req       = request.CustomerBatchRequest{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] GetBatchCustomersAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] GetBatchCustomersAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] GetBatchCustomersAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	results, err := u.userService.GetBatchCustomers(ctx, req.IDUsers)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetBatchCustomersAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	for _, result := range results {
		respBatch = append(respBatch, response.CustomerBatchResponse{
			ID:      result.ID,
			Name:    result.Name,
			Email:   result.Email,
			Phone:   result.Phone,
			Address: result.Address,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respBatch))
}

// DeleteCustomerAdmin implements UserHandlerInterface.
func (u *userHandler) DeleteCustomerAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] DeleteCustomerAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[UserHandler] DeleteCustomerAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[UserHandler] DeleteCustomerAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := u.userService.DeleteCustomer(ctx, id); err != nil {
		c.Logger().Errorf("[UserHandler] DeleteCustomerAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusNoContent, response.ResponseSuccess(nil))
}

// UpdateCustomerAdmin implements UserHandlerInterface.
func (u *userHandler) UpdateCustomerAdmin(c echo.Context) error {
	var (
		req = request.CustomerRequest{}
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] UpdateCustomerAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[UserHandler] UpdateCustomerAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[UserHandler] UpdateCustomerAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] UpdateCustomerAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] UpdateCustomerAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.UserEntity{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
		Photo:    req.Photo,
		RoleID:   2,
	}

	if err := u.userService.UpdateCustomer(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[UserHandler] UpdateCustomerAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrEmailAlreadyExists.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreateCustomerAdmin implements UserHandlerInterface.
func (u *userHandler) CreateCustomerAdmin(c echo.Context) error {
	var (
		req = request.CustomerRequest{}
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] CreateCustomerAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] CreateCustomerAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] CreateCustomerAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   2,
		Address:  req.Address,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	customerId, err := u.userService.CreateCustomer(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler] CreateCustomerAdmin: %v", err)
		switch err.Error() {
		case utils.ErrEmailAlreadyExists.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		case utils.ErrEmailNotVerified.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respCustomerId := map[string]int64{
		"customer_id": customerId,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respCustomerId))
}

// GetCustomerByIdAdmin implements UserHandlerInterface.
func (u *userHandler) GetCustomerByIdAdmin(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		respUser = response.CustomerResponse{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] GetCustomerByIdAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[UserHandler] GetCustomerByIdAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetCustomerByIdAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	result, err := u.userService.GetCustomerById(ctx, id)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetCustomerByIdAdmin: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respUser = response.CustomerResponse{
		ID:      result.ID,
		Name:    result.Name,
		Email:   result.Email,
		RoleID:  result.RoleID,
		Phone:   result.Phone,
		Lat:     result.Lat,
		Lng:     result.Lng,
		Address: result.Address,
		Photo:   result.Photo,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respUser))
}

// GetCustomersAllAdmin implements UserHandlerInterface.
func (u *userHandler) GetCustomersAllAdmin(c echo.Context) error {
	var (
		ctx      = c.Request().Context()
		respUser = []response.CustomerResponseList{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] GetCustomersAllAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	search := c.QueryParam("search")

	orderBy := c.QueryParam("order_by")
	if orderBy == "" {
		orderBy = "created_at"
	}

	orderType := c.QueryParam("order_type")
	if orderType != "asc" && orderType != "desc" {
		orderType = "desc"
	}

	page, err := conv.ParseInt64QueryParam(c, "page", 1)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetCustomersAllAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 5)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetCustomersAllAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.QueryStringEntity{
		Search:    search,
		OrderBy:   orderBy,
		OrderType: orderType,
		Page:      page,
		Limit:     limit,
	}

	results, countData, totalPages, err := u.userService.GetCustomersAll(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetCustomersAllAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	for _, val := range results {
		respUser = append(respUser, response.CustomerResponseList{
			ID:    val.ID,
			Name:  val.Name,
			Email: val.Email,
			Phone: val.Phone,
			Photo: val.Photo,
		})
	}

	pagination := response.Pagination{
		Page:       page,
		TotalCount: countData,
		PerPage:    limit,
		TotalPage:  totalPages,
	}

	return c.JSON(http.StatusOK, response.ResponseWithPaginationsSuccess(respUser, pagination))
}

// UpdateProfile implements UserHandlerInterface.
func (u *userHandler) UpdateProfile(c echo.Context) error {
	var (
		req         = request.UpdateDataRequest{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
		resp        = response.UpdateProfileResponse{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] UpdateProfile: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[UserHandler] UpdateProfile: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	userId := jwtUserData.UserID

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] UpdateProfile: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] UpdateProfile: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	latString := conv.LatLngToString(req.Lat)
	lngString := conv.LatLngToString(req.Lng)

	reqEntity := entity.UserEntity{
		ID:      userId,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Lat:     latString,
		Lng:     lngString,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	token, roleName, err := u.userService.UpdateProfile(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler] UpdateProfile: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrEmailAlreadyExists.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		case utils.ErrEmailNotVerified.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	resp = response.UpdateProfileResponse{
		ID:          userId,
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Photo:       req.Photo,
		Lat:         latString,
		Lng:         lngString,
		Role:        roleName,
		AccessToken: token,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(resp))
}

// GetProfileById implements UserHandlerInterface.
func (u *userHandler) GetProfileById(c echo.Context) error {
	var (
		respProfile = response.ProfileResponse{}
		ctx         = c.Request().Context()
		jwtUserData = entity.JwtUserData{}
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UserHandler] GetProfileById: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[UserHandler] GetProfileById: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	userId := jwtUserData.UserID

	result, err := u.userService.GetProfileById(ctx, userId)
	if err != nil {
		c.Logger().Errorf("[UserHandler] GetProfileById: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respProfile = response.ProfileResponse{
		ID:       result.ID,
		Name:     result.Name,
		Email:    result.Email,
		RoleName: result.RoleName,
		Phone:    result.Phone,
		Lat:      result.Lat,
		Lng:      result.Lng,
		Photo:    result.Photo,
		Address:  result.Address,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respProfile))
}

// UpdatePassword implements UserHandlerInterface.
func (u *userHandler) UpdatePassword(c echo.Context) error {
	var (
		req = request.UpdatePasswordRequest{}
		ctx = c.Request().Context()
	)

	tokenString := c.QueryParam("token")
	if tokenString == "" {
		c.Logger().Errorf("[UserHandler] UpdatePassword: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] UpdatePassword: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] UpdatePassword: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	if req.NewPassword != req.ConfirmPassword {
		c.Logger().Errorf("[UserHandler] UpdatePassword: %v", "New Password and Confirm Password do not match")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("New Password and Confirm Password do not match"))
	}

	reqEntity := entity.UserEntity{
		Password: req.NewPassword,
		Token:    tokenString,
	}

	if err := u.userService.UpdatePassword(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[UserHandler] UpdatePassword: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrTokenInvalid.Error():
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		case utils.ErrTokenExpired.Error():
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// ActivateAccount implements UserHandlerInterface.
func (u *userHandler) ActivateAccount(c echo.Context) error {
	var (
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	tokenString := c.QueryParam("token")
	if tokenString == "" {
		c.Logger().Errorf("[UserHandler] ActivateAccount: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	user, err := u.userService.ActivateAccount(ctx, tokenString)
	if err != nil {
		c.Logger().Errorf("[UserHandler] ActivateAccount: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrTokenInvalid.Error():
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		case utils.ErrTokenExpired.Error():
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respSignIn = response.SignInResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		AccessToken: user.Token,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respSignIn))
}

// ForgotPassword implements UserHandlerInterface.
func (u *userHandler) ForgotPassword(c echo.Context) error {
	var (
		req = request.ForgotPasswordRequest{}
		ctx = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] ForgotPassword: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] ForgotPassword: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.UserEntity{
		Email: req.Email,
		Token: uuid.New().String(),
	}

	if err := u.userService.ForgotPassword(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[UserHandler] ForgotPassword: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrEmailNotVerified.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// SignUp implements UserHandlerInterface.
func (u *userHandler) SignUp(c echo.Context) error {
	var (
		req = request.SignUpRequest{}
		ctx = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] SignUp: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] SignUp: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	if req.Password != req.PasswordConfirmation {
		c.Logger().Errorf("[UserHandler] SignUp: %v", "Password and Confirm Password do not match")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("Password and Confirm Password do not match"))
	}

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Token:    uuid.New().String(),
		RoleID:   2,
	}

	userId, err := u.userService.SignUp(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler] SignUp: %v", err)
		switch err.Error() {
		case utils.ErrEmailAlreadyExists.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		case utils.ErrEmailNotVerified.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respUserId := map[string]int64{
		"user_id": userId,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respUserId))
}

// SignIn implements UserHandlerInterface.
func (u *userHandler) SignIn(c echo.Context) error {
	var (
		req        = request.SignInRequest{}
		respSignIn = response.SignInResponse{}
		ctx        = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler] SignIn: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler] SignIn: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := u.userService.SignIn(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler] SignIn: %v", err)
		switch err.Error() {
		case utils.ErrDataNotFound.Error():
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.ErrLoginInvalid.Error():
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		case utils.ErrEmailNotVerified.Error():
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
		}
	}

	respSignIn = response.SignInResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Phone:       user.Phone,
		Lat:         user.Lat,
		Lng:         user.Lng,
		Role:        user.RoleName,
		Photo:       user.Photo,
		AccessToken: token,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respSignIn))
}
