package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	authHandler "github.com/pubestpubest/pos-backend/feature/auth/delivery"
	authRepository "github.com/pubestpubest/pos-backend/feature/auth/repository"
	authUsecase "github.com/pubestpubest/pos-backend/feature/auth/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func AuthRoutes(v1 *gin.RouterGroup) {
	authRepository := authRepository.NewAuthRepository(database.DB)
	authUsecase := authUsecase.NewAuthUsecase(authRepository)
	authHandler := authHandler.NewAuthHandler(authUsecase)

	authRoutes := v1.Group("/auth")
	{
		// Public routes
		authRoutes.POST("/login", authHandler.Login)

		// Protected routes
		protected := authRoutes.Group("")
		protected.Use(middlewares.AuthMiddleware())
		{
			protected.POST("/logout", authHandler.Logout)
			protected.POST("/change-password", authHandler.ChangePassword)
			protected.GET("/me", authHandler.GetMe)
		}
	}
}
