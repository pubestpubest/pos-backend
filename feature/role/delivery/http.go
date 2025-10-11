package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type roleHandler struct {
	roleUsecase domain.RoleUsecase
}

func NewRoleHandler(roleUsecase domain.RoleUsecase) *roleHandler {
	return &roleHandler{roleUsecase: roleUsecase}
}

func (h *roleHandler) GetAllRoles(c *gin.Context) {
	roles, err := h.roleUsecase.GetAllRoles()
	if err != nil {
		err = errors.Wrap(err, "[RoleHandler.GetAllRoles]: Error getting roles")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, roles)
}

func (h *roleHandler) GetRoleWithPermissions(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := h.roleUsecase.GetRoleWithPermissions(id)
	if err != nil {
		err = errors.Wrap(err, "[RoleHandler.GetRoleWithPermissions]: Error getting role")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, role)
}
