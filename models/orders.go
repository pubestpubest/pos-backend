package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	TableID      *uuid.UUID `gorm:"type:uuid;column:table_id"`
	OpenedBy     *uuid.UUID `gorm:"type:uuid;column:opened_by;comment:nullable if customer-originated is allowed"`
	Source       *string    `gorm:"type:varchar;column:source;comment:staff, customer"`
	Status       *string    `gorm:"type:varchar;column:status;comment:open, paid, void"`
	SubtotalBaht *int64     `gorm:"column:subtotal_baht"`
	DiscountBaht *int64     `gorm:"column:discount_baht"`
	TotalBaht    *int64     `gorm:"column:total_baht"`
	Note         *string    `gorm:"type:text;column:note"`
	CreatedAt    time.Time  `gorm:"type:timestamp;default:now();column:created_at"`
	ClosedAt     *time.Time `gorm:"type:timestamp;column:closed_at"`

	Table    *DiningTable `gorm:"foreignKey:TableID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
	Opener   *User        `gorm:"foreignKey:OpenedBy;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
	Items    []OrderItem  `gorm:"foreignKey:OrderID"`
	Payments []Payment    `gorm:"foreignKey:OrderID"`
}
