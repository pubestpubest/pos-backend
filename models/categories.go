package models

import "github.com/google/uuid"

type Category struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name         *string   `gorm:"type:varchar;uniqueIndex;column:name"`
	DisplayOrder *int      `gorm:"column:display_order"`

	MenuItems []MenuItem `gorm:"foreignKey:CategoryID"`
}
