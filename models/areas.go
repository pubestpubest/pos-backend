package models

import "github.com/google/uuid"

type Area struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name *string   `gorm:"type:varchar;uniqueIndex;column:name"`

	Tables []DiningTable `gorm:"foreignKey:AreaID"`
}
