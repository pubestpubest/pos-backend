package domain

import (
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/response"
)

// Permission domain - manages access permissions (mostly read-only, permissions are seeded)
type PermissionUsecase interface {
	GetAllPermissions() ([]*response.PermissionResponse, error)
}

type PermissionRepository interface {
	GetAllPermissions() ([]*models.Permission, error)
}
