package models

type Role struct {
	ID   int    `gorm:"primaryKey;autoIncrement;column:id"`
	Name string `gorm:"type:varchar;unique;not null;column:name;comment:cashier, waiter, kitchen, manager, owner"`

	// Associations
	Permissions []Permission `gorm:"many2many:role_permissions;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Users       []User       `gorm:"many2many:user_roles"`
}
