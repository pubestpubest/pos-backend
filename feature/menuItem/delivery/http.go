package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type menuItemHandler struct {
	menuItemUsecase domain.MenuItemUsecase
	minioClient     *minio.Client
}

func NewMenuItemHandler(menuItemUsecase domain.MenuItemUsecase, minioClient *minio.Client) *menuItemHandler {
	return &menuItemHandler{
		menuItemUsecase: menuItemUsecase,
		minioClient:     minioClient,
	}
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
	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get form values
	name := c.PostForm("name")
	sku := c.PostForm("sku")
	priceBahtStr := c.PostForm("price_baht")
	categoryIDStr := c.PostForm("category_id")
	activeStr := c.PostForm("active")

	// Validate required fields
	if name == "" || sku == "" || priceBahtStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, SKU, and price_baht are required"})
		return
	}

	// Parse price
	priceBaht, err := strconv.ParseInt(priceBahtStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price_baht value"})
		return
	}

	// Build request
	req := &request.MenuItemRequest{
		Name:      name,
		SKU:       sku,
		PriceBaht: priceBaht,
	}

	// Parse category ID if provided
	if categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
			return
		}
		req.CategoryID = &categoryID
	}

	// Parse active status
	if activeStr != "" {
		active := activeStr == "true"
		req.Active = &active
	}

	// Handle image upload if provided
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		imageURL, err := utils.UploadImageToMinio(h.minioClient, file)
		if err != nil {
			err = errors.Wrap(err, "[MenuItemHandler.CreateMenuItem]: Error uploading image")
			log.Warn(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
			return
		}
		req.ImageURL = &imageURL
	}

	menuItem, err := h.menuItemUsecase.CreateMenuItem(req)
	if err != nil {
		// If image was uploaded, clean it up
		if req.ImageURL != nil {
			_ = utils.DeleteImageFromMinio(h.minioClient, *req.ImageURL)
		}
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

	// Get existing menu item to get old image URL
	existingMenuItem, err := h.menuItemUsecase.GetMenuItemByID(id)
	if err != nil {
		err = errors.Wrap(err, "[MenuItemHandler.UpdateMenuItem]: Menu item not found")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}

	// Parse multipart form
	if err := c.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB max
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse form data"})
		return
	}

	// Get form values
	name := c.PostForm("name")
	sku := c.PostForm("sku")
	priceBahtStr := c.PostForm("price_baht")
	categoryIDStr := c.PostForm("category_id")
	activeStr := c.PostForm("active")

	// Validate required fields
	if name == "" || sku == "" || priceBahtStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, SKU, and price_baht are required"})
		return
	}

	// Parse price
	priceBaht, err := strconv.ParseInt(priceBahtStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price_baht value"})
		return
	}

	// Build request
	req := &request.MenuItemRequest{
		Name:      name,
		SKU:       sku,
		PriceBaht: priceBaht,
	}

	// Parse category ID if provided
	if categoryIDStr != "" {
		categoryID, err := uuid.Parse(categoryIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category_id"})
			return
		}
		req.CategoryID = &categoryID
	}

	// Parse active status
	if activeStr != "" {
		active := activeStr == "true"
		req.Active = &active
	}

	// Keep existing image URL by default
	oldImageURL := existingMenuItem.ImageURL
	req.ImageURL = &oldImageURL

	// Handle new image upload if provided
	file, err := c.FormFile("image")
	if err == nil && file != nil {
		newImageURL, err := utils.UploadImageToMinio(h.minioClient, file)
		if err != nil {
			err = errors.Wrap(err, "[MenuItemHandler.UpdateMenuItem]: Error uploading image")
			log.Warn(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
			return
		}
		req.ImageURL = &newImageURL

		// Delete old image if it exists and is different
		if oldImageURL != "" && oldImageURL != newImageURL {
			if err := utils.DeleteImageFromMinio(h.minioClient, oldImageURL); err != nil {
				log.Warnf("[MenuItemHandler.UpdateMenuItem]: Failed to delete old image: %v", err)
			}
		}
	}

	menuItem, err := h.menuItemUsecase.UpdateMenuItem(id, req)
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

	// Soft delete - just mark as deleted in database
	// Image is kept in storage in case item needs to be restored
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
