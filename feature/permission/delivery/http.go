package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type permissionHandler struct {
	permissionUsecase domain.PermissionUsecase
}

func NewPermissionHandler(permissionUsecase domain.PermissionUsecase) *permissionHandler {
	return &permissionHandler{permissionUsecase: permissionUsecase}
}

func (h *permissionHandler) GetAllPermissions(c *gin.Context) {
	permissions, err := h.permissionUsecase.GetAllPermissions()
	if err != nil {
		err = errors.Wrap(err, "[PermissionHandler.GetAllPermissions]: Error getting permissions")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, permissions)
}
