package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// MenuItem domain - manages menu items and their available modifiers
type MenuItemUsecase interface {
	GetAllMenuItems() ([]*response.MenuItemResponse, error)
	GetMenuItemByID(id uuid.UUID) (*response.MenuItemResponse, error)
	CreateMenuItem(req *request.MenuItemRequest) (*response.MenuItemResponse, error)
	UpdateMenuItem(id uuid.UUID, req *request.MenuItemRequest) (*response.MenuItemResponse, error)
	DeleteMenuItem(id uuid.UUID) error
	GetAvailableModifiers() ([]*response.ModifierResponse, error)
}

type MenuItemRepository interface {
	GetAllMenuItems() ([]*models.MenuItem, error)
	GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error)
	CreateMenuItem(menuItem *models.MenuItem) error
	UpdateMenuItem(menuItem *models.MenuItem) error
	DeleteMenuItem(id uuid.UUID) error
	GetAllModifiers() ([]*models.Modifier, error)
}
