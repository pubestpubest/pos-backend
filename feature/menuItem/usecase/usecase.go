package usecase

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type menuItemUsecase struct {
	menuItemRepository domain.MenuItemRepository
}

func NewMenuItemUsecase(menuItemRepository domain.MenuItemRepository) domain.MenuItemUsecase {
	return &menuItemUsecase{menuItemRepository: menuItemRepository}
}

func (u *menuItemUsecase) GetAllMenuItems() ([]*response.MenuItemResponse, error) {
	menuItems, err := u.menuItemRepository.GetAllMenuItems()
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetAllMenuItems]: Error getting menu items")
	}

	menuItemResponses := make([]*response.MenuItemResponse, len(menuItems))
	for i, menuItem := range menuItems {
		menuItemResponses[i] = &response.MenuItemResponse{
			ID:        menuItem.ID,
			Name:      utils.DerefString(menuItem.Name),
			PriceBaht: utils.DerefInt64(menuItem.PriceBaht),
			Active:    utils.DerefBool(menuItem.Active),
			ImageURL:  utils.DerefString(menuItem.ImageURL),
			Category:  response.CategoryResponse{ID: utils.DerefUUID(menuItem.CategoryID), Name: utils.DerefString(menuItem.Category.Name), DisplayOrder: utils.DerefInt(menuItem.Category.DisplayOrder)},
		}
	}

	return menuItemResponses, nil
}

func (u *menuItemUsecase) GetMenuItemByID(id uuid.UUID) (*response.MenuItemResponse, error) {
	menuItem, err := u.menuItemRepository.GetMenuItemByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetMenuItemByID]: Error getting menu item")
	}

	return &response.MenuItemResponse{
		ID:        menuItem.ID,
		Name:      utils.DerefString(menuItem.Name),
		PriceBaht: utils.DerefInt64(menuItem.PriceBaht),
		Active:    utils.DerefBool(menuItem.Active),
		ImageURL:  utils.DerefString(menuItem.ImageURL),
		Category:  response.CategoryResponse{ID: utils.DerefUUID(menuItem.CategoryID), Name: utils.DerefString(menuItem.Category.Name), DisplayOrder: utils.DerefInt(menuItem.Category.DisplayOrder)},
	}, nil
}

func (u *menuItemUsecase) CreateMenuItem(req *request.MenuItemRequest) (*response.MenuItemResponse, error) {
	active := true
	if req.Active != nil {
		active = *req.Active
	}

	menuItem := &models.MenuItem{
		CategoryID: req.CategoryID,
		Name:       &req.Name,
		SKU:        &req.SKU,
		PriceBaht:  &req.PriceBaht,
		Active:     &active,
		ImageURL:   req.ImageURL,
	}

	if err := u.menuItemRepository.CreateMenuItem(menuItem); err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.CreateMenuItem]: Error creating menu item")
	}

	return &response.MenuItemResponse{
		ID:        menuItem.ID,
		Name:      req.Name,
		PriceBaht: req.PriceBaht,
		Active:    active,
		ImageURL:  utils.DerefString(req.ImageURL),
		Category:  response.CategoryResponse{ID: utils.DerefUUID(req.CategoryID), Name: "", DisplayOrder: 0},
	}, nil
}

func (u *menuItemUsecase) UpdateMenuItem(id uuid.UUID, req *request.MenuItemRequest) (*response.MenuItemResponse, error) {
	// Get existing menu item
	menuItem, err := u.menuItemRepository.GetMenuItemByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.UpdateMenuItem]: Menu item not found")
	}

	// Update fields
	menuItem.CategoryID = req.CategoryID
	menuItem.Name = &req.Name
	menuItem.SKU = &req.SKU
	menuItem.PriceBaht = &req.PriceBaht
	menuItem.ImageURL = req.ImageURL
	if req.Active != nil {
		menuItem.Active = req.Active
	}

	if err := u.menuItemRepository.UpdateMenuItem(menuItem); err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.UpdateMenuItem]: Error updating menu item")
	}

	return &response.MenuItemResponse{
		ID:        menuItem.ID,
		Name:      utils.DerefString(menuItem.Name),
		PriceBaht: utils.DerefInt64(menuItem.PriceBaht),
		Active:    utils.DerefBool(menuItem.Active),
		ImageURL:  utils.DerefString(menuItem.ImageURL),
		Category:  response.CategoryResponse{ID: utils.DerefUUID(menuItem.CategoryID), Name: utils.DerefString(menuItem.Category.Name), DisplayOrder: utils.DerefInt(menuItem.Category.DisplayOrder)},
	}, nil
}

func (u *menuItemUsecase) DeleteMenuItem(id uuid.UUID) error {
	// Check if menu item exists
	_, err := u.menuItemRepository.GetMenuItemByID(id)
	if err != nil {
		return errors.Wrap(err, "[MenuItemUsecase.DeleteMenuItem]: Menu item not found")
	}

	if err := u.menuItemRepository.DeleteMenuItem(id); err != nil {
		return errors.Wrap(err, "[MenuItemUsecase.DeleteMenuItem]: Error deleting menu item")
	}

	return nil
}

func (u *menuItemUsecase) GetAvailableModifiers() ([]*response.ModifierResponse, error) {
	modifiers, err := u.menuItemRepository.GetAllModifiers()
	if err != nil {
		return nil, errors.Wrap(err, "[MenuItemUsecase.GetAvailableModifiers]: Error getting modifiers")
	}

	modifierResponses := make([]*response.ModifierResponse, len(modifiers))
	for i, modifier := range modifiers {
		modifierResponses[i] = &response.ModifierResponse{
			ID:             modifier.ID,
			Name:           utils.DerefString(modifier.Name),
			PriceDeltaBaht: utils.DerefInt64(modifier.PriceDeltaBaht),
			Note:           utils.DerefString(modifier.Note),
		}
	}

	return modifierResponses, nil
}
