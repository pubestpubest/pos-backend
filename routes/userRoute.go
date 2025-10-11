package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	userHandler "github.com/pubestpubest/pos-backend/feature/user/delivery"
	userRepository "github.com/pubestpubest/pos-backend/feature/user/repository"
	userUsecase "github.com/pubestpubest/pos-backend/feature/user/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func UserRoutes(v1 *gin.RouterGroup) {
	userRepository := userRepository.NewUserRepository(database.DB)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userHandler := userHandler.NewUserHandler(userUsecase)

	userRoutes := v1.Group("/users")
	userRoutes.Use(middlewares.AuthMiddleware())
	{
		userRoutes.GET("", userHandler.GetAllUsers)
		userRoutes.GET("/:id", userHandler.GetUserByID)
		userRoutes.POST("", userHandler.CreateUser)
		userRoutes.PUT("/:id", userHandler.UpdateUser)
		userRoutes.POST("/:id/roles", userHandler.AssignRole)
	}
}
