package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey;column:id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;column:user_id;index"`
	Token     string    `gorm:"type:varchar;unique;not null;column:token;index"`
	ExpiresAt time.Time `gorm:"type:timestamp;not null;column:expires_at"`
	CreatedAt time.Time `gorm:"type:timestamp;default:now();column:created_at"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
