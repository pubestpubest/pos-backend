package models

import "github.com/google/uuid"

type DiningTable struct {
	ID     uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	AreaID *uuid.UUID `gorm:"type:uuid;column:area_id"`
	Name   *string    `gorm:"type:varchar;column:name"`
	Seats  *int       `gorm:"column:seats"`
	Status *string    `gorm:"type:varchar;column:status;comment:free, occupied, needs_pay"`
	QRSlug *string    `gorm:"type:varchar;unique;column:qr_slug;comment:unguessable slug used in QR URLs"`

	Area   *Area   `gorm:"foreignKey:AreaID;references:ID;constraint:OnUpdate:SET NULL,OnDelete:SET NULL"`
	Orders []Order `gorm:"foreignKey:TableID"`
}
