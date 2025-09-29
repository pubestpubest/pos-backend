package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
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
	table, err := h.tableUsecase.GetAllTables()
	if err != nil {
		err = errors.Wrap(err, "[TableHandler.GetAllTables]: Error getting user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}
	c.JSON(http.StatusOK, table)
}
