package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user-service/config"
	"user-service/internal/adapter/handler"
	"user-service/internal/adapter/repository"
	"user-service/internal/adapter/repository/cache"
	"user-service/internal/adapter/storage"
	"user-service/internal/core/service"
	"user-service/utils/logger"
	"user-service/utils/validator"

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

	storageHandler := storage.NewSupabase(cfg)

	txManager := repository.NewGormTransactionManager(db.DB)

	userRepo := repository.NewUserRepository(db.DB, customLogger.Logger())
	tokenRepo := repository.NewVerificationTokenRepository(db.DB, customLogger.Logger())
	roleRepo := repository.NewRoleRepository(db.DB, customLogger.Logger())
	outboxEventRepo := repository.NewOutboxEventRepository(db.DB, customLogger.Logger())

	cacheUser := cache.NewUserCache(redisClient, userRepo, customLogger.Logger())

	jwtService := service.NewJwtService(cfg)
	roleService := service.NewRoleService(roleRepo, redisClient, txManager, customLogger.Logger())
	userService := service.NewUserService(userRepo, cfg, jwtService, tokenRepo, outboxEventRepo, roleService, cacheUser, txManager, customLogger.Logger())

	handler.NewUserHandler(e, userService, cfg, jwtService, redisClient)
	handler.NewUploadImageStorageHandler(e, cfg, jwtService, storageHandler, redisClient)
	handler.NewRoleHandler(e, roleService, cfg, jwtService, redisClient)

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

	customLogger.Logger().Infof("[RunServer-4] shutting down server on 5 seconds...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
