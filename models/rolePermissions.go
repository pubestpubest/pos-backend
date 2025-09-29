package models

type RolePermission struct {
	RoleID       int `gorm:"not null;primaryKey;column:role_id"`
	PermissionID int `gorm:"not null;primaryKey;column:permission_id"`

	Role       *Role       `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission *Permission `gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
