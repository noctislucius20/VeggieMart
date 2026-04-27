package seeds

import (
	"context"
	"log"
	"strings"
	"user-service/internal/core/domain/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedRole(ctx context.Context, db *gorm.DB) {
	var (
		permissions    []model.Permission
		rolePermission []model.RolePermission
		roles          = []model.Role{
			{
				Name: "Super Admin",
			},
			{
				Name: "Customer",
			},
		}
	)

	if err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("id", "resource", "scope").
			Find(&permissions).Error; err != nil {
			return err
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(&roles, 100).Error; err != nil {
			return err
		}

		if err := tx.Select("id", "name").
			Find(&roles).Error; err != nil {
			return err
		}

		for _, r := range roles {
			for _, p := range permissions {
				if strings.ToLower(p.Scope) == "all" && r.Name == "Super Admin" {
					rolePermission = append(rolePermission, model.RolePermission{
						RoleID:       r.ID,
						PermissionID: p.ID,
					})
				}
				if strings.ToLower(p.Scope) == "own" {
					rolePermission = append(rolePermission, model.RolePermission{
						RoleID:       r.ID,
						PermissionID: p.ID,
					})
				}
			}
		}

		if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
			CreateInBatches(&rolePermission, 100).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
		return
	}

	log.Println("All roles created")
}
