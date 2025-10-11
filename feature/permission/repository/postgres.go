package repository

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) domain.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) GetAllPermissions() ([]*models.Permission, error) {
	var permissions []*models.Permission
	if err := r.db.Order("code ASC").Find(&permissions).Error; err != nil {
		return nil, errors.Wrap(err, "[PermissionRepository.GetAllPermissions]: Error querying database")
	}
	return permissions, nil
}
