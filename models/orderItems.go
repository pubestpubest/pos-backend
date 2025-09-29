package models

import "github.com/google/uuid"

type OrderItem struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	OrderID       uuid.UUID `gorm:"type:uuid;not null;column:order_id"`
	MenuItemID    uuid.UUID `gorm:"type:uuid;not null;column:menu_item_id"`
	Quantity      int       `gorm:"column:quantity"`
	UnitPriceBaht int64     `gorm:"column:unit_price_baht"`
	LineTotalBaht int64     `gorm:"column:line_total_baht"`
	Note          *string   `gorm:"type:text;column:note"`

	Order    *Order    `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	MenuItem *MenuItem `gorm:"foreignKey:MenuItemID;references:ID;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT"`

	Modifiers []OrderItemModifier `gorm:"foreignKey:OrderItemID"`
}
