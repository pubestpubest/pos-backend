package models

import "github.com/google/uuid"

type OrderItemModifier struct {
	OrderItemID    uuid.UUID `gorm:"type:uuid;not null;primaryKey;column:order_item_id"`
	ModifierID     uuid.UUID `gorm:"type:uuid;not null;primaryKey;column:modifier_id"`
	PriceDeltaBaht *int64    `gorm:"column:price_delta_baht"`

	OrderItem *OrderItem `gorm:"foreignKey:OrderItemID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Modifier  *Modifier  `gorm:"foreignKey:ModifierID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
