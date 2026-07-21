package main

import (
	"api-gateway/config"
	"api-gateway/handler"
	"api-gateway/middleware/cors"
	"api-gateway/middleware/jwt"
	"api-gateway/middleware/logger"
	ratelimit "api-gateway/middleware/rate-limit"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

func main() {
	serviceCtx, serviceCancel := context.WithCancel(context.Background())
	defer serviceCancel()

	logrusLogger := logrus.New()
	logrusLogger.SetFormatter(&logrus.JSONFormatter{})

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient, err := config.NewRedisClient(serviceCtx)
	if err != nil {
		logrusLogger.Fatal(err)
		return
	}

	e := echo.New()

	redisRateLimit := ratelimit.NewRedisRateLimiter(redisClient)
	jwtMiddleware := jwt.NewMiddlewareJWT(redisClient)

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(logger.MiddlewareLog())
	e.Use(cors.MiddlewareCORS())
	e.Use(redisRateLimit.MiddlewareRateLimiter())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "OK"})
	})

	api := e.Group("/api")

	handler.RegisterAllRoutes(api, jwtMiddleware.CheckToken())

	go func() {
		var port string

		if os.Getenv("APP_PORT") != "" {
			port = os.Getenv("APP_PORT")
		}

		err = e.Start(":" + port)
		if err != nil {
			logrusLogger.Fatal(err)
			return
		}

		logrusLogger.Infof("API Gateway starting on port: %s", port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	logrusLogger.Infof("shutting down API gateway on 5 seconds...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	e.Shutdown(ctx)
}
