package usecase

import (
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
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
		err = errors.Wrap(err, "[MenuItemUsecase.GetAllMenuItems]: Error getting menu item")
		return nil, err
	}

	menuItemResponses := make([]*response.MenuItemResponse, len(menuItems))
	for i, menuItem := range menuItems {
		menuItemResponses[i] = &response.MenuItemResponse{
			ID:         menuItem.ID,
			Name:       utils.DerefString(menuItem.Name),
			PriceBaht:  utils.DerefInt64(menuItem.PriceBaht),
			Active:     utils.DerefBool(menuItem.Active),
			ImageURL:   utils.DerefString(menuItem.ImageURL),
			CategoryID: utils.DerefUUID(menuItem.CategoryID),
		}
	}

	return menuItemResponses, nil
}
