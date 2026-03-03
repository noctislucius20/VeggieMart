package handler

import (
	"net/http"
	"strconv"
	"time"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/request"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/service"
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

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.GET("/roles", roleHandler.GetRolesAllAdmin)
	adminGroup.GET("/roles/:id", roleHandler.GetRoleByIdAdmin)
	adminGroup.POST("/roles", roleHandler.CreateRoleAdmin)
	adminGroup.PUT("/roles/:id", roleHandler.UpdateRoleAdmin)
	adminGroup.DELETE("/roles/:id", roleHandler.DeleteRoleAdmin)

	return roleHandler
}

// CreateRoleAdmin implements RoleHandlerInterface.
func (r *roleHandler) CreateRoleAdmin(c echo.Context) error {
	var (
		req = request.RoleRequest{}
		ctx = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[RoleHandler-1] CreateRoleAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[RoleHandler-2] CreateRoleAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[RoleHandler-3] CreateRoleAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	roleEntity := entity.RoleEntity{
		Name: req.Name,
	}

	roleId, err := r.roleService.CreateRoleAdmin(ctx, roleEntity)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-4] CreateRoleAdmin: %v", err)
		if err.Error() == utils.DATA_ALREADY_EXISTS {
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[RoleHandler-1] DeleteRoleAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[RoleHandler-2] DeleteRoleAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	roleId, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-3] DeleteRoleAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	err = r.roleService.DeleteRoleAdmin(ctx, roleId)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-4] DeleteRoleAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		if err.Error() == utils.DATA_STILL_IN_USED {
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}

// GetRolesAllAdmin implements RoleHandlerInterface.
func (r *roleHandler) GetRolesAllAdmin(c echo.Context) error {
	var (
		respRole = []response.RoleResponse{}
		ctx      = c.Request().Context()
	)

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[RoleHandler-1] GetRolesAllAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	search := c.QueryParam("search")

	roles, err := r.roleService.GetRolesAllAdmin(ctx, search)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-2] GetRolesAllAdmin: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[RoleHandler-1] GetRoleByIdAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[RoleHandler-2] GetRoleByIdAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-3] GetRoleByIdAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	role, err := r.roleService.GetRoleByIdAdmin(ctx, id)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-4] GetRoleByIdAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
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

	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[RoleHandler-1] UpdateRoleAdmin: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed("data token invalid"))
	}

	idParam := c.Param("id")
	if idParam == "" {
		c.Logger().Errorf("[RoleHandler-2] UpdateRoleAdmin: %v", "id required")
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("id required"))
	}

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-3] UpdateRoleAdmin: %v", "id invalid")
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed("id invalid"))
	}

	if err := c.Bind(&req); err != nil {
		c.Logger().Errorf("[RoleHandler-4] UpdateRoleAdmin: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	if err := c.Validate(&req); err != nil {
		c.Logger().Errorf("[RoleHandler-5] UpdateRoleAdmin: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	reqEntity := entity.RoleEntity{
		ID:   id,
		Name: req.Name,
	}

	err = r.roleService.UpdateRoleAdmin(ctx, reqEntity)
	if err != nil {
		c.Logger().Errorf("[RoleHandler-6] UpdateRoleAdmin: %v", err)
		if err.Error() == utils.DATA_NOT_FOUND {
			return c.JSON(http.StatusNotFound, response.ResponseFailed(err.Error()))
		}
		if err.Error() == utils.DATA_ALREADY_EXISTS {
			return c.JSON(http.StatusConflict, response.ResponseFailed(err.Error()))
		}
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(utils.INTERNAL_SERVER_ERROR))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(nil))
}
