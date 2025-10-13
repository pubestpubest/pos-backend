package delivery

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type authHandler struct {
	authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *authHandler {
	return &authHandler{authUsecase: authUsecase}
}

func (h *authHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	authResponse, err := h.authUsecase.Login(&req)
	if err != nil {
		err = errors.Wrap(err, "[AuthHandler.Login]: Error logging in")
		log.Warn(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.StandardError(err)})
		return
	}

	// Set HTTP-only cookie with the token
	c.SetCookie("token", authResponse.Token, 24*60*60, "/", "", false, true) // 24 hours, HTTP-only, secure=false for development

	// Return response without token
	response := gin.H{
		"user":        authResponse.User,
		"expires_at":  authResponse.ExpiresAt,
		"permissions": authResponse.Permissions,
	}

	c.JSON(http.StatusOK, response)
}

func (h *authHandler) Logout(c *gin.Context) {
	// Get token from cookie
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authentication cookie not found"})
		return
	}

	if err := h.authUsecase.Logout(token); err != nil {
		err = errors.Wrap(err, "[AuthHandler.Logout]: Error logging out")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}

	// Clear the cookie
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *authHandler) ChangePassword(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.authUsecase.ChangePassword(userID.(uuid.UUID), &req); err != nil {
		err = errors.Wrap(err, "[AuthHandler.ChangePassword]: Error changing password")
		log.Warn(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.StandardError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *authHandler) GetMe(c *gin.Context) {
	// Get user from context (set by auth middleware)
	user, exists := c.Get("user")
	if !exists {
		log.Warn("[AuthHandler.GetMe]: User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fmt.Println("Get user passed")
	// Get user permissions
	userID := c.GetString("userID")
	fmt.Println("userID: ", userID)
	uuid, err := uuid.Parse(userID)
	if err != nil {
		log.Warn("[AuthHandler.GetMe]: Invalid user ID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	permissions, err := h.authUsecase.GetUserPermissions(uuid)
	if err != nil {
		err = errors.Wrap(err, "[AuthHandler.GetMe]: Error getting permissions")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":        user,
		"permissions": permissions,
	})
}
