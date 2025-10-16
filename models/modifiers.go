package models

import "github.com/google/uuid"

type Modifier struct {
	ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	CategoryID     *uuid.UUID `gorm:"type:uuid;column:category_id"`
	Name           *string    `gorm:"type:varchar;uniqueIndex;column:name"`
	PriceDeltaBaht *int64     `gorm:"column:price_delta_baht;default:0"`
	Note           *string    `gorm:"type:text;column:note"`

	Category *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
}
