package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"product-service/config"
	"product-service/internal/adapter"
	"product-service/internal/adapter/handler/response"
	"product-service/internal/adapter/storage"
	"product-service/internal/core/service"
	"product-service/utils"
	"product-service/utils/logger"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type UploadImageInterface interface {
	UploadImage(c echo.Context) error
}

type uploadImageStruct struct {
	storageHandler storage.SupabaseInterface
}

// UploadImage implements UploadImageInterface.
func (u *uploadImageStruct) UploadImage(c echo.Context) error {
	user := c.Get("user").(string)
	if user == "" {
		c.Logger().Errorf("[UploadImageHandler-1] UploadImage: %v", "data token not found")
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.TOKEN_INVALID))
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.Logger().Errorf("[UploadImageHandler-1] UploadImage: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("[UploadImageHandler-2] UploadImage: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	defer src.Close()

	fileBuffer := new(bytes.Buffer)
	_, err = io.Copy(fileBuffer, src)
	if err != nil {
		c.Logger().Errorf("[UploadImageHandler-3] UploadImage: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), path.Ext(file.Filename))

	uploadPath := fmt.Sprintf("public/uploads/%s", newFileName)
	url, err := u.storageHandler.UploadFile(uploadPath, fileBuffer)
	if err != nil {
		c.Logger().Errorf("[UploadImageHandler-4] UploadImage: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	respImgUrl := map[string]string{"image_url": url}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respImgUrl))

}

func NewUploadImageStorageHandler(e *echo.Echo, cfg *config.Config, jwtService service.JwtServiceInterface, storageHandler storage.SupabaseInterface, redisClient *redis.Client) UploadImageInterface {
	uploadImageHandler := &uploadImageStruct{
		storageHandler: storageHandler,
	}

	e.Use(middleware.Recover())
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Timeout: 10 * time.Second,
	}))

	mid := adapter.NewMiddlewareAdapter(cfg, logger.NewLogger().Logger(), jwtService, redisClient)

	adminGroup := e.Group("/admin", mid.CheckToken())
	adminGroup.POST("/image-upload", uploadImageHandler.UploadImage)

	return uploadImageHandler
}
