package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
	"time"
	"user-service/config"
	"user-service/internal/adapter"
	"user-service/internal/adapter/handler/response"
	"user-service/internal/adapter/storage"
	"user-service/internal/core/service"
	middlewareGateway "user-service/internal/middleware"
	"user-service/utils"
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

	authUploadGroup := userGroup.Group("/profile", mid.CheckToken(), mid.RequiredPermission(authPermission...))
	authUploadGroup.POST("/image-upload", res.UploadImage)

	return res
}

// UploadImage implements UploadImageInterface.
func (u *uploadImageStruct) UploadImage(c echo.Context) error {
	user, ok := c.Get("user").(string)
	if !ok || user == "" {
		c.Logger().Errorf("[UploadImage] UploadImage: %v", utils.ErrTokenInvalid.Error())
		return c.JSON(http.StatusUnauthorized, response.ResponseFailed(utils.ErrTokenInvalid.Error()))
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.Logger().Errorf("[UploadImage] UploadImage: %v", err)
		return c.JSON(http.StatusUnprocessableEntity, response.ResponseFailed(err.Error()))
	}

	// Validasi ukuran file (max 1 MB)
	const maxFileSize = 1 * 1024 * 1024 // 1 MB
	if file.Size > maxFileSize {
		c.Logger().Errorf("[UploadImage] File size exceeds limit: %d bytes", file.Size)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("file size exceeds maximum limit of 1 MB"))
	}

	// Validasi tipe file (hanya jpg/png)
	ext := strings.ToLower(path.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.Logger().Errorf("[UploadImage] Invalid file extension: %s", ext)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("only JPG and PNG files are allowed"))
	}

	// Validasi content-type
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		c.Logger().Errorf("[UploadImage] Invalid content type: %s", contentType)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed("only JPG and PNG files are allowed"))
	}

	src, err := file.Open()
	if err != nil {
		c.Logger().Errorf("[UploadImage] UploadImage: %v", err)
		return c.JSON(http.StatusBadRequest, response.ResponseFailed(err.Error()))
	}

	defer src.Close()

	fileBuffer := new(bytes.Buffer)
	_, err = io.Copy(fileBuffer, src)
	if err != nil {
		c.Logger().Errorf("[UploadImage] UploadImage: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	newFileName := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), path.Ext(file.Filename))

	uploadPath := fmt.Sprintf("public/uploads/%s", newFileName)
	url, err := u.storageHandler.UploadFile(uploadPath, fileBuffer)
	if err != nil {
		c.Logger().Errorf("[UploadImage] UploadImage: %v", err)
		return c.JSON(http.StatusInternalServerError, response.ResponseFailed(err.Error()))
	}

	respImgUrl := map[string]string{"image_url": url}

	return c.JSON(http.StatusOK, response.ResponseSuccess(respImgUrl))

}
