package repository

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type menuItemRepository struct {
	db *gorm.DB
}

func NewMenuItemRepository(db *gorm.DB) domain.MenuItemRepository {
	return &menuItemRepository{db: db}
}

func (r *menuItemRepository) GetAllMenuItems() ([]*models.MenuItem, error) {
	var menuItemsList []*models.MenuItem
	if err := r.db.Find(&menuItemsList).Error; err != nil {
		err = errors.Wrap(err, "[MenuItemRepository.GetAllMenuItems]: Error getting menu items")
		return nil, err
	}

	return menuItemsList, nil
}
