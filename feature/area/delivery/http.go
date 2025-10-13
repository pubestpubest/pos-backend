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

type areaHandler struct {
	areaUsecase domain.AreaUsecase
}

func NewAreaHandler(areaUsecase domain.AreaUsecase) *areaHandler {
	return &areaHandler{areaUsecase: areaUsecase}
}

func (h *areaHandler) GetAllAreas(c *gin.Context) {
	areas, err := h.areaUsecase.GetAllAreas()
	if err != nil {
		err = errors.Wrap(err, "[AreaHandler.GetAllAreas]: Error getting areas")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, areas)
}

func (h *areaHandler) GetAreaByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
		return
	}

	area, err := h.areaUsecase.GetAreaByID(id)
	if err != nil {
		err = errors.Wrap(err, "[AreaHandler.GetAreaByID]: Error getting area")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, area)
}

func (h *areaHandler) CreateArea(c *gin.Context) {
	var req request.AreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	area, err := h.areaUsecase.CreateArea(&req)
	if err != nil {
		err = errors.Wrap(err, "[AreaHandler.CreateArea]: Error creating area")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusCreated, area)
}

func (h *areaHandler) UpdateArea(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
		return
	}

	var req request.AreaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	area, err := h.areaUsecase.UpdateArea(id, &req)
	if err != nil {
		err = errors.Wrap(err, "[AreaHandler.UpdateArea]: Error updating area")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, area)
}

func (h *areaHandler) DeleteArea(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid area ID"})
		return
	}

	if err := h.areaUsecase.DeleteArea(id); err != nil {
		err = errors.Wrap(err, "[AreaHandler.DeleteArea]: Error deleting area")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Area deleted successfully"})
}

func (h *areaHandler) GetAreasWithTables(c *gin.Context) {
	areas, err := h.areaUsecase.GetAreasWithTables()
	if err != nil {
		err = errors.Wrap(err, "[AreaHandler.GetAreasWithTables]: Error getting areas with tables")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, areas)
}
