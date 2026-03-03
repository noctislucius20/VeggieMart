package app

import (
	"context"
	"order-service/config"
	"order-service/internal/adapter/handler"
	httpclient "order-service/internal/adapter/http_client"
	"order-service/internal/adapter/repository"
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
		customLogger.Logger().Fatalf("[RunServer-1] %v", err.Error())
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

	esClient, err := cfg.NewElasticsearchClient()
	if err != nil {
		customLogger.Logger().Fatalf("[RunServer-3] %v", err.Error())
		return
	}

	httpClient := httpclient.NewHttpClient(cfg)

	orderRepo := repository.NewOrderRepository(customLogger.Logger())
	outboxRepo := repository.NewOutboxEventRepository(customLogger.Logger())
	elasticRepo := repository.NewElasticRepository(esClient, customLogger.Logger())

	httpService := service.NewHttpService(cfg, httpClient)
	orderService := service.NewOrderService(orderRepo, outboxRepo, elasticRepo, httpService, db.DB, customLogger.Logger())

	handler.NewOrderHandler(orderService, e, cfg)

	e.GET("/api/check", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	go func() {
		if cfg.App.AppPort == "" {
			cfg.App.AppPort = os.Getenv("APP_PORT")
		}

		err = e.Start(":" + cfg.App.AppPort)
		if err != nil {
			customLogger.Logger().Fatalf("[RunServer-4] %v", err.Error())
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
