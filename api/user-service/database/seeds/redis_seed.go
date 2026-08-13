package seeds

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func RedisPermissionSeed(ctx context.Context, redisClient *redis.Client, db *gorm.DB) {
	var (
		roles []model.Role
		pipe  = redisClient.Pipeline()
	)

	if err := db.Select("id", "name").
		Preload("Permissions", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "resource", "action", "scope")
		}).
		Find(&roles).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
		return
	}

	for _, r := range roles {
		roleEntity := entity.RoleEntity{
			ID:   r.ID,
			Name: r.Name,
		}

		for _, p := range r.Permissions {
			roleEntity.Permissions = append(roleEntity.Permissions, entity.PermissionEntity{
				ID:       p.ID,
				Resource: p.Resource,
				Action:   p.Action,
				Scope:    p.Scope,
			})
		}

		jsonData, err := json.Marshal(roleEntity)
		if err != nil {
			log.Fatalf("%s: %v", err.Error(), err)
			return
		}

		if err := pipe.Set(ctx, fmt.Sprintf("role:id:%d", roleEntity.ID), jsonData, 0).Err(); err != nil {
			log.Fatalf("%s: %v", err.Error(), err)
			return
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
		return
	}

	log.Println("Redis Permission Seed Completed")
}
