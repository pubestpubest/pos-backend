package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type menuItemHandler struct {
	menuItemUsecase domain.MenuItemUsecase
}

func NewMenuItemHandler(menuItemUsecase domain.MenuItemUsecase) *menuItemHandler {
	return &menuItemHandler{menuItemUsecase: menuItemUsecase}
}

func (h *menuItemHandler) GetAllMenuItems(c *gin.Context) {
	menuItem, err := h.menuItemUsecase.GetAllMenuItems()
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.GetAllMenuItems]: Error getting user")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
	}
	c.JSON(http.StatusOK, menuItem)
}
