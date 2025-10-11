package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	categoryHandler "github.com/pubestpubest/pos-backend/feature/category/delivery"
	categoryRepository "github.com/pubestpubest/pos-backend/feature/category/repository"
	categoryUsecase "github.com/pubestpubest/pos-backend/feature/category/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func CategoryRoutes(v1 *gin.RouterGroup) {
	categoryRepository := categoryRepository.NewCategoryRepository(database.DB)
	categoryUsecase := categoryUsecase.NewCategoryUsecase(categoryRepository)
	categoryHandler := categoryHandler.NewCategoryHandler(categoryUsecase)

	categoryRoutes := v1.Group("/categories")
	categoryRoutes.Use(middlewares.AuthMiddleware())
	{
		categoryRoutes.GET("", categoryHandler.GetAllCategories)
		categoryRoutes.GET("/:id", categoryHandler.GetCategoryByID)
		categoryRoutes.POST("", categoryHandler.CreateCategory)
		categoryRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		categoryRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
