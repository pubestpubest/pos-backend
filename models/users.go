package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	Username     string    `gorm:"type:varchar;unique;not null;column:username"`
	PasswordHash string    `gorm:"type:text;not null;column:password_hash"`
	FullName     *string   `gorm:"type:varchar;column:full_name"`
	Email        *string   `gorm:"type:varchar;column:email"`
	Phone        *string   `gorm:"type:varchar;column:phone"`
	Status       *string   `gorm:"type:varchar;column:status;comment:active, locked"`
	CreatedAt    time.Time `gorm:"type:timestamp;default:now();column:created_at"`
	UpdatedAt    time.Time `gorm:"type:timestamp;default:now();column:updated_at"`

	// Associations
	Roles []Role `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}
