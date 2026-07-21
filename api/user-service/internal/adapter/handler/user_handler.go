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
	CreateUserAccount(c echo.Context) error
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
	userGroup.POST("/signup", userHandler.CreateUserAccount)
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
	userGroup.GET("/customers", userHandler.GetCustomersAllAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	userGroup.POST("/customers/batch", userHandler.GetBatchCustomersAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	userGroup.GET("/customers/:id", userHandler.GetCustomerByIdAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	userGroup.POST("/customers", userHandler.CreateCustomerAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	userGroup.PUT("/customers/:id", userHandler.UpdateCustomerAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	userGroup.DELETE("/customers/:id", userHandler.DeleteCustomerAdmin, mid.CheckToken(), mid.RequiredPermission(adminPermission...))

	// authGroup := e.Group("/auth", mid.CheckToken())
	userGroup.GET("/profile", userHandler.GetProfileById, mid.CheckToken(), mid.RequiredPermission(authPermission...))
	userGroup.PUT("/profile", userHandler.UpdateProfile, mid.CheckToken(), mid.RequiredPermission(authPermission...))

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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] GetBatchCustomersAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler-2] GetBatchCustomersAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-3] GetBatchCustomersAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	results, err := u.userService.GetBatchCustomers(ctx, req.IDUsers)
	if err != nil {
		c.Logger().Errorf("[UserHandler-4] GetBatchCustomersAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] DeleteCustomerAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[UserHandler-2] DeleteCustomerAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_INVALID))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[UserHandler-3] DeleteCustomerAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	if err := u.userService.DeleteCustomer(ctx, id); err != nil {
		c.Logger().Errorf("[UserHandler-4] DeleteCustomerAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// UpdateCustomerAdmin implements UserHandlerInterface.
func (u *userHandler) UpdateCustomerAdmin(c echo.Context) error {
	var (
		req = request.CustomerRequest{}
		ctx = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] UpdateCustomerAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[UserHandler-2] UpdateCustomerAdmin: %s", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_INVALID))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[UserHandler-3] UpdateCustomerAdmin: %s", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler-4] UpdateCustomerAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-5] UpdateCustomerAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	latString := conv.LatLngToString(req.Lat)
	lngString := conv.LatLngToString(req.Lng)

	reqEntity := entity.UserEntity{
		ID:       id,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Phone:    req.Phone,
		Address:  req.Address,
		Lat:      latString,
		Lng:      lngString,
		Photo:    req.Photo,
		RoleID:   2,
	}

	if err := u.userService.UpdateCustomer(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[UserHandler-6] UpdateCustomerAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.EMAIL_ALREADY_EXISTS:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] CreateCustomerAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler-2] CreateCustomerAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-3] CreateCustomerAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	latString := conv.LatLngToString(req.Lat)
	lngString := conv.LatLngToString(req.Lng)

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		RoleID:   2,
		Address:  req.Address,
		Lat:      latString,
		Lng:      lngString,
		Phone:    req.Phone,
		Photo:    req.Photo,
	}

	customerId, err := u.userService.CreateCustomer(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler-5] CreateCustomerAdmin: %v", err)
		switch err.Error() {
		case utils.EMAIL_ALREADY_EXISTS:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		case utils.EMAIL_NOT_VERIFIED:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] GetCustomerByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[UserHandler-2] GetCustomerByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ID_INVALID))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[UserHandler-3] GetCustomerByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ID_INVALID))
	}

	result, err := u.userService.GetCustomerById(ctx, id)
	if err != nil {
		c.Logger().Errorf("[UserHandler-4] GetCustomerByIdAdmin: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] GetCustomersAllAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
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
		c.Logger().Errorf("[UserHandler-2] GetCustomersAllAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	limit, err := conv.ParseInt64QueryParam(c, "limit", 10)
	if err != nil {
		c.Logger().Errorf("[UserHandler-3] GetCustomersAllAdmin: %v", err)
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
		c.Logger().Errorf("[UserHandler-4] GetCustomersAllAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] UpdateProfile: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[UserHandler-2] UpdateProfile: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	userId := jwtUserData.UserID

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler-3] UpdateProfile: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-4] UpdateProfile: %v", err)
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
		c.Logger().Errorf("[UserHandler-5] UpdateProfile: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.EMAIL_ALREADY_EXISTS:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		case utils.EMAIL_NOT_VERIFIED:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UserHandler-1] GetProfileById: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := json.Unmarshal([]byte(user), &jwtUserData); err != nil {
		c.Logger().Errorf("[UserHandler-2] GetProfileById: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	userId := jwtUserData.UserID

	result, err := u.userService.GetProfileById(ctx, userId)
	if err != nil {
		c.Logger().Errorf("[UserHandler-3] GetProfileById: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
		c.Logger().Errorf("[UserHandler-1] UpdatePassword: %v", "missing or invalid token")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler-2] UpdatePassword: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-3] UpdatePassword: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	if req.NewPassword != req.ConfirmPassword {
		c.Logger().Errorf("[UserHandler-4] UpdatePassword: %v", "New Password and Confirm Password do not match")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("New Password and Confirm Password do not match"))
	}

	reqEntity := entity.UserEntity{
		Password: req.NewPassword,
		Token:    tokenString,
	}

	if err := u.userService.UpdatePassword(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[UserHandler-5] UpdatePassword: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.TOKEN_INVALID:
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		case utils.TOKEN_EXPIRED:
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
		c.Logger().Errorf("[UserHandler-1] ActivateAccount: %v", "missing or invalid token")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	user, err := u.userService.ActivateAccount(ctx, tokenString)
	if err != nil {
		c.Logger().Errorf("[UserHandler-2] ActivateAccount: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.TOKEN_INVALID:
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		case utils.TOKEN_EXPIRED:
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
		c.Logger().Errorf("[UserHandler-1] ForgotPassword: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-2] ForgotPassword: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.UserEntity{
		Email: req.Email,
		Token: uuid.New().String(),
	}

	if err := u.userService.ForgotPassword(ctx, reqEntity); err != nil {
		c.Logger().Errorf("[UserHandler-3] ForgotPassword: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.EMAIL_NOT_VERIFIED:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
		}
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// CreateUserAccount implements UserHandlerInterface.
func (u *userHandler) CreateUserAccount(c echo.Context) error {
	var (
		req = request.SignUpRequest{}
		ctx = c.Request().Context()
	)

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[UserHandler-1] CreateUserAccount: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-2] CreateUserAccount: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	if req.Password != req.PasswordConfirmation {
		c.Logger().Errorf("[UserHandler-3] CreateUserAccount: %v", "Password and Confirm Password do not match")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("Password and Confirm Password do not match"))
	}

	reqEntity := entity.UserEntity{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Token:    uuid.New().String(),
		RoleID:   2,
	}

	userId, err := u.userService.CreateUserAccount(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler-4] CreateUserAccount: %v", err)
		switch err.Error() {
		case utils.EMAIL_ALREADY_EXISTS:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		case utils.EMAIL_NOT_VERIFIED:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
		c.Logger().Errorf("[UserHandler-1] SignIn: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[UserHandler-2] SignIn: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.UserEntity{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := u.userService.SignIn(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[UserHandler-3] SignIn: %v", err)
		switch err.Error() {
		case utils.DATA_NOT_FOUND:
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		case utils.LOGIN_INVALID:
			return c.JSON(http.StatusUnauthorized, response.ResponseFailed(err.Error()))
		case utils.EMAIL_NOT_VERIFIED:
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		default:
			return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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
