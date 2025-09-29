package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	OrderID     uuid.UUID `gorm:"type:uuid;not null;column:order_id"`
	Method      *string   `gorm:"type:varchar;column:method;comment:cash, card, promptpay"`
	AmountBaht  int64     `gorm:"column:amount_baht"`
	Currency    *string   `gorm:"type:varchar(3);default:THB;column:currency"`
	Provider    *string   `gorm:"type:varchar;column:provider"`
	ProviderRef *string   `gorm:"type:varchar;column:provider_ref"`
	Status      *string   `gorm:"type:varchar;column:status;comment:succeeded, pending, failed"`
	CreatedAt   time.Time `gorm:"type:timestamp;default:now();column:created_at"`

	Order *Order `gorm:"foreignKey:OrderID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
