package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuItem struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	CategoryID *uuid.UUID     `gorm:"type:uuid;column:category_id"`
	Name       *string        `gorm:"type:varchar;column:name"`
	SKU        *string        `gorm:"type:varchar;unique;column:sku"`
	PriceBaht  *int64         `gorm:"column:price_baht"`
	Active     *bool          `gorm:"column:active;default:true"`
	ImageURL   *string        `gorm:"type:text;column:image_url"`
	DeletedAt  gorm.DeletedAt `gorm:"index;column:deleted_at"`

	Category *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
}
