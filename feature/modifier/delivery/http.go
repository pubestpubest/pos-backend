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

type modifierHandler struct {
	modifierUsecase domain.ModifierUsecase
}

func NewModifierHandler(modifierUsecase domain.ModifierUsecase) *modifierHandler {
	return &modifierHandler{modifierUsecase: modifierUsecase}
}

func (h *modifierHandler) GetAllModifiers(c *gin.Context) {
	modifiers, err := h.modifierUsecase.GetAllModifiers()
	if err != nil {
		err = errors.Wrap(err, "[ModifierHandler.GetAllModifiers]: Error getting modifiers")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, modifiers)
}

func (h *modifierHandler) GetModifierByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid modifier ID"})
		return
	}

	modifier, err := h.modifierUsecase.GetModifierByID(id)
	if err != nil {
		err = errors.Wrap(err, "[ModifierHandler.GetModifierByID]: Error getting modifier")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, modifier)
}

func (h *modifierHandler) CreateModifier(c *gin.Context) {
	var req request.ModifierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	modifier, err := h.modifierUsecase.CreateModifier(&req)
	if err != nil {
		err = errors.Wrap(err, "[ModifierHandler.CreateModifier]: Error creating modifier")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusCreated, modifier)
}

func (h *modifierHandler) UpdateModifier(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid modifier ID"})
		return
	}

	var req request.ModifierRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	modifier, err := h.modifierUsecase.UpdateModifier(id, &req)
	if err != nil {
		err = errors.Wrap(err, "[ModifierHandler.UpdateModifier]: Error updating modifier")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, modifier)
}

func (h *modifierHandler) DeleteModifier(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid modifier ID"})
		return
	}

	if err := h.modifierUsecase.DeleteModifier(id); err != nil {
		err = errors.Wrap(err, "[ModifierHandler.DeleteModifier]: Error deleting modifier")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Modifier deleted successfully"})
}

func (h *modifierHandler) GetModifiersByCategoryID(c *gin.Context) {
	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	modifiers, err := h.modifierUsecase.GetModifiersByCategoryID(categoryID)
	if err != nil {
		err = errors.Wrap(err, "[ModifierHandler.GetModifiersByCategoryID]: Error getting modifiers")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, modifiers)
}
