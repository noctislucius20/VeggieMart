package app

import (
	"context"
	"os"
	"os/signal"
	"product-service/config"
	"product-service/internal/adapter/handler"
	"product-service/internal/adapter/repository"
	"product-service/internal/adapter/storage"
	"product-service/internal/core/service"
	"product-service/utils/logger"
	"product-service/utils/validator"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RunServer() {
	serviceCtx, serviceCancel := context.WithCancel(context.Background())

	customValidator := validator.NewValidator()
	customLogger := logger.NewLogger()

	cfg := config.NewConfig()
	db, err := cfg.ConnectionPostgres(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-1] %v", err)
		return
	}

	redisClient, err := cfg.NewRedisClient(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-2] %v", err)
		return
	}

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogMethod:     true,
		LogStatus:     true,
		LogValuesFunc: customLogger.RequestLogger,
	}))

	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)

	e.Validator = customValidator
	e.Logger = customLogger.Logger()

	esClient, err := cfg.NewElasticSearchClient()
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-2] %v", err.Error())
		return
	}

	storageHandler := storage.NewSupabase(cfg, customLogger.Logger())

	txManager := repository.NewGormTransactionManager(db.DB)

	categoryRepo := repository.NewCategoryRepository(db.DB, customLogger.Logger())
	productRepo := repository.NewProductRepository(db.DB, esClient, customLogger.Logger())
	outboxEventRepo := repository.NewOutboxEventRepository(db.DB, customLogger.Logger())

	jwtService := service.NewJwtService(cfg)
	categoryService := service.NewCategoryService(categoryRepo, redisClient, txManager, customLogger.Logger())
	productService := service.NewProductService(cfg, productRepo, redisClient, txManager, categoryService, outboxEventRepo, customLogger.Logger())

	handler.NewCategoryHandler(e, categoryService, cfg, jwtService, redisClient)
	handler.NewProductHandler(e, cfg, productService, jwtService, redisClient)
	handler.NewUploadImageStorageHandler(e, cfg, jwtService, storageHandler, redisClient)

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err = e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer-3] %v", err.Error())
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	serviceCancel()

	customLogger.Logger().Infof("[RunServer-5] shutting down server on 5 seconds...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
