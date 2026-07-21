package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"time"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/adapter/storage"
	"user-service/internal/core/service"
	middlewareGateway "user-service/internal/middleware"
	"user-service/utils/logger"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type UploadImageInterface interface {
	UploadImage(c echo.Context) error
}

type uploadImageStruct struct {
	storageHandler storage.SupabaseInterface
}

// UploadImage implements UploadImageInterface.
func (u *uploadImageStruct) UploadImage(c echo.Context) error {
	file, err := c.FormFile("photo")
	if err != nil {
		c.Logger().Errorf("[UploadImage-1] UploadImage: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("[UploadImage-2] UploadImage: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	defer src.Close()

	fileBuffer := new(bytes.Buffer)
	_, err = io.Copy(fileBuffer, src)
	if err != nil {
		c.Logger().Errorf("[UploadImage-3] UploadImage: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), path.Ext(file.Filename))

	uploadPath := fmt.Sprintf("public/uploads/%s", newFileName)
	url, err := u.storageHandler.UploadFile(uploadPath, fileBuffer)
	if err != nil {
		c.Logger().Errorf("[UploadImage-4] UploadImage: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	return c.JSON(http.StatusOK, response.ResponseSuccess(map[string]string{"image_url": url}))

}

func NewUploadImageStorageHandler(e *echo.Echo, cfg *config.Config, jwtService service.JwtServiceInterface, storageHandler storage.SupabaseInterface, redisClient *redis.Client) UploadImageInterface {
	res := &uploadImageStruct{
		storageHandler: storageHandler,
	}

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	authPermission := []string{
		"users:read:own",
		"users:write:own",
		"users:update:own",
		"users:delete:own",
	}

	userGroup := e.Group("/users")
	userGroup.Use(middlewareGateway.GatewayValidationMiddleware(cfg))

	userGroup.POST("/profile/image-upload", res.UploadImage, mid.CheckToken(), mid.RequiredPermission(authPermission...))

	return res
}
