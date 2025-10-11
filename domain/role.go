package domain

import (
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/response"
)

// Role domain - manages user roles (mostly read-only, roles are seeded)
type RoleUsecase interface {
	GetAllRoles() ([]*response.RoleResponse, error)
	GetRoleWithPermissions(id int) (*response.RoleResponse, error)
}

type RoleRepository interface {
	GetAllRoles() ([]*models.Role, error)
	GetRoleWithPermissions(id int) (*models.Role, error)
}
