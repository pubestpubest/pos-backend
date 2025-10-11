package repository

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) domain.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) GetAllRoles() ([]*models.Role, error) {
	var roles []*models.Role
	if err := r.db.Order("name ASC").Find(&roles).Error; err != nil {
		return nil, errors.Wrap(err, "[RoleRepository.GetAllRoles]: Error querying database")
	}
	return roles, nil
}

func (r *roleRepository) GetRoleWithPermissions(id int) (*models.Role, error) {
	var role models.Role
	if err := r.db.Preload("Permissions").Where("id = ?", id).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[RoleRepository.GetRoleWithPermissions]: Role not found")
		}
		return nil, errors.Wrap(err, "[RoleRepository.GetRoleWithPermissions]: Error querying database")
	}
	return &role, nil
}
