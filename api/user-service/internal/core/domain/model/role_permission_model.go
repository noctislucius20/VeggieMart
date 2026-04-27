package model

type RolePermission struct {
	PermissionID int64 `gorm:"primaryKey"`
	RoleID       int64 `gorm:"primaryKey"`

	Permission Permission `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Role       Role       `gorm:"constraint:OnUpdate:CASCADE;"`
}

func (RolePermission) TableName() string {
	return "role_permission"
}
