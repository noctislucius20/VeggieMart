package permissionseed

import (
	"log"
	"user-service/database/seeds/permission_seed/data"
	"user-service/internal/core/domain/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func SeedPermission(db *gorm.DB) {
	permission := []model.Permission{}

	permission = append(permission, data.UserPermissions...)
	permission = append(permission, data.RolePermissions...)
	permission = append(permission, data.CategoryPermissions...)
	permission = append(permission, data.ProductPermissions...)
	permission = append(permission, data.OrderPermissions...)
	permission = append(permission, data.PaymentPermissions...)

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(&permission, 100).Error; err != nil {
		log.Fatalf("%s: %v", err.Error(), err)
		return
	}

	log.Printf("All permissions created")

}
