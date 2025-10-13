package response

import (
	"time"
)

type AuthResponse struct {
	User        UserResponse `json:"user"`
	Token       string       // Internal field for setting cookie, not serialized
	ExpiresAt   time.Time    `json:"expires_at"`
	Permissions []string     `json:"permissions"`
}

type PermissionCheckResponse struct {
	HasPermission bool `json:"has_permission"`
}
