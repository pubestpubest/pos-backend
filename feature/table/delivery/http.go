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

type tableHandler struct {
	tableUsecase domain.TableUsecase
}

func NewTableHandler(tableUsecase domain.TableUsecase) *tableHandler {
	return &tableHandler{tableUsecase: tableUsecase}
}

func (h *tableHandler) GetAllTables(c *gin.Context) {
	tables, err := h.tableUsecase.GetAllTables()
	if err != nil {
		err = errors.Wrap(err, "[TableHandler.GetAllTables]: Error getting tables")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, tables)
}

func (h *tableHandler) GetTableByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	}

	table, err := h.tableUsecase.GetTableByID(id)
	if err != nil {
		err = errors.Wrap(err, "[TableHandler.GetTableByID]: Error getting table")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, table)
}

func (h *tableHandler) UpdateTableStatus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	}

	var req request.UpdateTableStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.tableUsecase.UpdateTableStatus(id, req.Status); err != nil {
		err = errors.Wrap(err, "[TableHandler.UpdateTableStatus]: Error updating table status")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Table status updated successfully"})
}
