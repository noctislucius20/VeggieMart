package model

type UserRole struct {
	UserID int64 `gorm:"primaryKey"`
	RoleID int64 `gorm:"primaryKey"`

	User User `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
	Role Role `gorm:"constraint:OnUpdate:CASCADE;"`
}

func (UserRole) TableName() string {
	return "user_role"
}
