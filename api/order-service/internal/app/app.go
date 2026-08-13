package app

import (
	"context"
	"order-service/config"
	"order-service/internal/adapter/handler"
	httpclient "order-service/internal/adapter/http_client"
	"order-service/internal/adapter/repository"
	"order-service/internal/adapter/repository/cache"
	"order-service/internal/core/service"
	"order-service/utils/logger"
	"order-service/utils/validator"
	"os"
	"os/signal"
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
		customLogger.Logger().Fatalf("[RunServer] %v", err)
		return
	}

	redisClient, err := cfg.NewRedisClient(serviceCtx)
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer] %v", err)
		return
	}

	e := echo.New()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:        true,
		LogMethod:     true,
		LogStatus:     true,
		LogValuesFunc: customLogger.RequestLogger,
	}))

	en.RegisterDefaultTranslations(customValidator.Validator, customValidator.Translator)

	e.Validator = customValidator
	e.Logger = customLogger.Logger()

	esClient, err := cfg.NewElasticsearchClient()
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer] %v", err)
		return
	}

	httpClient := httpclient.NewHttpClient(cfg)

	txManager := repository.NewGormTransactionManager(db.DB)

	orderRepo := repository.NewOrderRepository(db.DB, customLogger.Logger())
	outboxRepo := repository.NewOutboxEventRepository(db.DB, customLogger.Logger())
	elasticRepo := repository.NewElasticRepository(esClient, customLogger.Logger())

	cacheOrder := cache.NewOrderCache(redisClient, orderRepo, customLogger.Logger())

	jwtService := service.NewJwtService(cfg)
	httpService := service.NewHttpService(cfg, httpClient)
	orderService := service.NewOrderService(cfg, orderRepo, outboxRepo, elasticRepo, cacheOrder, httpService, txManager, customLogger.Logger())

	handler.NewOrderHandler(e, cfg, orderService, jwtService, redisClient)

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err = e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer] %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	serviceCancel()

	customLogger.Logger().Infof("[RunServer] shutting down server on 5 seconds...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
