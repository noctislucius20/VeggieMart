package handler

import (
	"errors"
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
	"user-service/utils/logger"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type RoleHandlerInterface interface {
	GetRolesAllAdmin(c echo.Context) error
	GetRoleByIdAdmin(c echo.Context) error
	CreateRoleAdmin(c echo.Context) error
	UpdateRoleAdmin(c echo.Context) error
	DeleteRoleAdmin(c echo.Context) error
}

type roleHandler struct {
	roleService service.RoleServiceInterface
}

func NewRoleHandler(e *echo.Echo, roleService service.RoleServiceInterface, cfg *config.Config, jwtService service.JwtServiceInterface, redisClient *redis.Client) RoleHandlerInterface {
	roleHandler := &roleHandler{roleService: roleService}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	adminPermission := []string{
		"roles:read:all",
		"roles:write:all",
		"roles:update:all",
		"roles:delete:all",
	}

	// adminGroup := e.Group("/admin", mid.CheckToken())
	userGroup := e.Group("/users")
	userGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	adminRoleGroup := userGroup.Group("/roles", mid.CheckToken(), mid.RequiredPermission(adminPermission...))
	adminRoleGroup.GET("", roleHandler.GetRolesAllAdmin)
	adminRoleGroup.GET("/:id", roleHandler.GetRoleByIdAdmin)
	adminRoleGroup.POST("", roleHandler.CreateRoleAdmin)
	adminRoleGroup.PUT("/:id", roleHandler.UpdateRoleAdmin)
	adminRoleGroup.DELETE("/:id", roleHandler.DeleteRoleAdmin)

	return roleHandler
}

// CreateRoleAdmin implements RoleHandlerInterface.
func (r *roleHandler) CreateRoleAdmin(c echo.Context) error {
	var (
		req = request.RoleRequest{}
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[RoleHandler] CreateRoleAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[RoleHandler] CreateRoleAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[RoleHandler] CreateRoleAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	roleEntity := entity.RoleEntity{
		Name: req.Name,
	}

	roleId, err := r.roleService.CreateRoleAdmin(ctx, roleEntity)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] CreateRoleAdmin: %v", err)
		if errors.Is(err, utils.ErrDataAlreadyExists) {
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	respRoleId := map[string]int64{
		"role_id": roleId,
	}

	return c.JSON(http.StatusCreated, response.ResponseSuccess(respRoleId))
}

// DeleteRoleAdmin implements RoleHandlerInterface.
func (r *roleHandler) DeleteRoleAdmin(c echo.Context) error {
	var (
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[RoleHandler] DeleteRoleAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[RoleHandler] DeleteRoleAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	roleId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] DeleteRoleAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	err = r.roleService.DeleteRoleAdmin(ctx, roleId)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] DeleteRoleAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		if errors.Is(err, utils.ErrDataStillInUsed) {
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusNoContent, response.ResponseSuccess(nil))
}

// GetRolesAllAdmin implements RoleHandlerInterface.
func (r *roleHandler) GetRolesAllAdmin(c echo.Context) error {
	var (
		respRole = []response.RoleResponse{}
		ctx      = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[RoleHandler] GetRolesAllAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	search := c.QueryParam("search")

	roles, err := r.roleService.GetRolesAllAdmin(ctx, search)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] GetRolesAllAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	for _, role := range roles {
		respRole = append(respRole, response.RoleResponse{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respRole))
}

// GetRoleByIdAdmin implements RoleHandlerInterface.
func (r *roleHandler) GetRoleByIdAdmin(c echo.Context) error {
	var (
		respRole = response.RoleResponse{}
		ctx      = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[RoleHandler] GetRoleByIdAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[RoleHandler] GetRoleByIdAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] GetRoleByIdAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	role, err := r.roleService.GetRoleByIdAdmin(ctx, id)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] GetRoleByIdAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	respRole = response.RoleResponse{
		ID:   role.ID,
		Name: role.Name,
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respRole))
}

// UpdateRoleAdmin implements RoleHandlerInterface.
func (r *roleHandler) UpdateRoleAdmin(c echo.Context) error {
	var (
		req = request.RoleRequest{}
		ctx = c.Request().Context()
	)

	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[RoleHandler] UpdateRoleAdmin: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[RoleHandler] UpdateRoleAdmin: %v", utils.ErrIDRequired.Error())
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(utils.ErrIDRequired.Error()))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] UpdateRoleAdmin: %v", utils.ErrIDInvalid.Error())
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(utils.ErrIDInvalid.Error()))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[RoleHandler] UpdateRoleAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[RoleHandler] UpdateRoleAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.RoleEntity{
		ID:   id,
		Name: req.Name,
	}

	err = r.roleService.UpdateRoleAdmin(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[RoleHandler] UpdateRoleAdmin: %v", err)
		if errors.Is(err, utils.ErrDataNotFound) {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		if errors.Is(err, utils.ErrDataAlreadyExists) {
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.ErrInternalServerError.Error()))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}
