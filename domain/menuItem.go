package domain

import (
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/response"
)

type MenuItemUsecase interface {
	GetAllMenuItems() ([]*response.MenuItemResponse, error)
}

type MenuItemRepository interface {
	GetAllMenuItems() ([]*models.MenuItem, error)
}
