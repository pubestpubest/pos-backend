package models

import "github.com/google/uuid"

type Modifier struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Name           *string   `gorm:"type:varchar;uniqueIndex;column:name"`
	PriceDeltaBaht *int64    `gorm:"column:price_delta_baht;default:0"`
	Note           *string   `gorm:"type:text;column:note"`
}
