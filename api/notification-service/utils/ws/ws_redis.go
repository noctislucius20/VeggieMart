package ws

import (
	"context"
	"encoding/json"
	"notification-service/internal/core/domain/entity"
	"notification-service/internal/core/service"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

const (
	WS_REDIS_CHANNEL       = "ws:notifications"
	WS_ADMIN_REDIS_CHANNEL = "ws:admin_notifications"
)

func SubscribeWebSocketChannel(ctx context.Context, redisClient *redis.Client, notificationService service.NotificationServiceInterface, logger *log.Logger) {
	pubsub := redisClient.Subscribe(ctx, WS_REDIS_CHANNEL)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			var wsMsg entity.WsRedisEntity
			if err := json.Unmarshal([]byte(msg.Payload), &wsMsg); err != nil {
				logger.Errorf("[WsRedis] SubscribeWebSocketChannel-1: %v", err)
				continue
			}

			conn := GetWebSocketConn(wsMsg.ReceiverID)
			if conn == nil {
				logger.Errorf("[WsRedis] SubscribeWebSocketChannel-2: WebSocket connection not found for user %d", wsMsg.ReceiverID)
				continue
			}

			if err := conn.WriteJSON(msg); err != nil {
				logger.Errorf("[WsRedis] SubscribeWebSocketChannel-3: %v", err)
				RemoveWebSocketConn(wsMsg.ReceiverID)
				continue
			}
		}
	}
}

func SubscribeAdminWebSocketChannel(ctx context.Context, redisClient *redis.Client, logger *log.Logger) {
	pubsub := redisClient.Subscribe(ctx, WS_ADMIN_REDIS_CHANNEL)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-ch:
			if !ok {
				return
			}

			logger.Infof("[WsRedis] SubscribeAdminWebSocketChannel: received message on admin channel")
			BroadcastAdminMessage(msg)
		}
	}
}

func PublishWebSocketMessage(ctx context.Context, redisClient *redis.Client, wsMsg entity.WsRedisEntity, logger *log.Logger) error {
	data, err := json.Marshal(wsMsg)
	if err != nil {
		logger.Errorf("[WsRedis] PublishWebSocketMessage-1: %v", err)
		return err
	}

	if err := redisClient.Publish(ctx, WS_REDIS_CHANNEL, string(data)).Err(); err != nil {
		logger.Errorf("[WsRedis] PublishWebSocketMessage-2: %v", err)
		return err
	}

	return nil
}

func PublishAdminWebSocketMessage(ctx context.Context, redisClient *redis.Client, wsMsg entity.WsRedisEntity, logger *log.Logger) error {
	data, err := json.Marshal(wsMsg)
	if err != nil {
		logger.Errorf("[WsRedis] PublishAdminWebSocketMessage-1: %v", err)
		return err
	}

	if err := redisClient.Publish(ctx, WS_ADMIN_REDIS_CHANNEL, string(data)).Err(); err != nil {
		logger.Errorf("[WsRedis] PublishAdminWebSocketMessage-2: %v", err)
		return err
	}

	return nil
}
