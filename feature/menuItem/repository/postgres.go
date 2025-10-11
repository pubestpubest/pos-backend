package repository

import (
	"github.com/google/uuid"
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
	if err := r.db.Preload("Category").Order("name ASC").Find(&menuItemsList).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetAllMenuItems]: Error getting menu items")
	}
	return menuItemsList, nil
}

func (r *menuItemRepository) GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Preload("Category").Where("id = ?", id).First(&menuItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemByID]: Menu item not found")
		}
		return nil, errors.Wrap(err, "[MenuItemRepository.GetMenuItemByID]: Error querying database")
	}
	return &menuItem, nil
}

func (r *menuItemRepository) CreateMenuItem(menuItem *models.MenuItem) error {
	if err := r.db.Create(menuItem).Error; err != nil {
		return errors.Wrap(err, "[MenuItemRepository.CreateMenuItem]: Error creating menu item")
	}
	return nil
}

func (r *menuItemRepository) UpdateMenuItem(menuItem *models.MenuItem) error {
	if err := r.db.Save(menuItem).Error; err != nil {
		return errors.Wrap(err, "[MenuItemRepository.UpdateMenuItem]: Error updating menu item")
	}
	return nil
}

func (r *menuItemRepository) DeleteMenuItem(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.MenuItem{}).Error; err != nil {
		return errors.Wrap(err, "[MenuItemRepository.DeleteMenuItem]: Error deleting menu item")
	}
	return nil
}

func (r *menuItemRepository) GetAllModifiers() ([]*models.Modifier, error) {
	var modifiers []*models.Modifier
	if err := r.db.Order("name ASC").Find(&modifiers).Error; err != nil {
		return nil, errors.Wrap(err, "[MenuItemRepository.GetAllModifiers]: Error getting modifiers")
	}
	return modifiers, nil
}
