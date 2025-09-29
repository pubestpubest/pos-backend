package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/go-clean-arch-template/database"
	userHandler "github.com/pubestpubest/go-clean-arch-template/feature/user/delivery"
	userRepository "github.com/pubestpubest/go-clean-arch-template/feature/user/repository"
	userUsecase "github.com/pubestpubest/go-clean-arch-template/feature/user/usecase"
)

func UserRoutes(v1 *gin.RouterGroup) {

	userRepository := userRepository.NewUserRepository(database.DB)
	userUsecase := userUsecase.NewUserUsecase(userRepository)
	userHandler := userHandler.NewUserHandler(userUsecase)

	userRoutes := v1.Group("/users")

	userRoutes.GET("/:id", userHandler.GetUser)
}
