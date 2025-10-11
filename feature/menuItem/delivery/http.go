package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/request"
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
	menuItems, err := h.menuItemUsecase.GetAllMenuItems()
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.GetAllMenuItems]: Error getting menu items")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, menuItems)
}

func (h *menuItemHandler) GetMenuItemByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	menuItem, err := h.menuItemUsecase.GetMenuItemByID(id)
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.GetMenuItemByID]: Error getting menu item")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, menuItem)
}

func (h *menuItemHandler) CreateMenuItem(c *gin.Context) {
	var req request.MenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	menuItem, err := h.menuItemUsecase.CreateMenuItem(&req)
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.CreateMenuItem]: Error creating menu item")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusCreated, menuItem)
}

func (h *menuItemHandler) UpdateMenuItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	var req request.MenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	menuItem, err := h.menuItemUsecase.UpdateMenuItem(id, &req)
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.UpdateMenuItem]: Error updating menu item")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, menuItem)
}

func (h *menuItemHandler) DeleteMenuItem(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid menu item ID"})
		return
	}

	if err := h.menuItemUsecase.DeleteMenuItem(id); err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.DeleteMenuItem]: Error deleting menu item")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}

func (h *menuItemHandler) GetAvailableModifiers(c *gin.Context) {
	modifiers, err := h.menuItemUsecase.GetAvailableModifiers()
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.GetAvailableModifiers]: Error getting modifiers")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, modifiers)
}
