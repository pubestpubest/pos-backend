package delivery

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/go-clean-arch-template/domain"
	"github.com/pubestpubest/go-clean-arch-template/utils"
	log "github.com/sirupsen/logrus"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *userHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		err = errors.Wrap(err, "[UserHandler.GetUser]: Error parsing id")
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}
	user, err := h.userUsecase.GetUser(uint32(idUint))
	if err != nil {
		err = errors.Wrap(err, "[UserHandler.GetUser]: Error getting user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		log.Warn(err)
		return
	}
	c.JSON(http.StatusOK, user)
}
