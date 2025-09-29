package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type userHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(userUsecase domain.UserUsecase) *userHandler {
	return &userHandler{userUsecase: userUsecase}
}

func (h *userHandler) GetAllUsers(c *gin.Context) {
	user, err := h.userUsecase.GetAllUsers()
	if err != nil {
		err = errors.Wrap(err, "[UserHandler.GetAllUsers]: Error getting user")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
	}
	c.JSON(http.StatusOK, user)
}
