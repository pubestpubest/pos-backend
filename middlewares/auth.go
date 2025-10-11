package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/database"
	authRepository "github.com/pubestpubest/pos-backend/feature/auth/repository"
	authUsecase "github.com/pubestpubest/pos-backend/feature/auth/usecase"
	"github.com/pubestpubest/pos-backend/response"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token (remove "Bearer " prefix)
		token := authHeader
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = authHeader[7:]
		}

		// Validate token
		authRepo := authRepository.NewAuthRepository(database.DB)
		authUc := authUsecase.NewAuthUsecase(authRepo)

		user, err := authUc.GetUserByToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set user in context
		c.Set("userID", user.ID)
		c.Set("user", response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Email:    user.Email,
			Phone:    user.Phone,
			Status:   user.Status,
		})

		c.Next()
	}
}

func RequirePermission(permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		authRepo := authRepository.NewAuthRepository(database.DB)
		authUc := authUsecase.NewAuthUsecase(authRepo)

		hasPermission, err := authUc.VerifyPermission(userID.(uuid.UUID), permissionCode)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
