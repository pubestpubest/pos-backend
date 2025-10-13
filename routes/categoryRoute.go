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

	// Public routes for customers
	categoryPublicRoutes := v1.Group("/categories")
	{
		categoryPublicRoutes.GET("", categoryHandler.GetAllCategories)
		categoryPublicRoutes.GET("/:id", categoryHandler.GetCategoryByID)
	}

	// Protected routes for staff/admin
	categoryProtectedRoutes := v1.Group("/categories")
	categoryProtectedRoutes.Use(middlewares.AuthMiddleware())
	{
		categoryProtectedRoutes.POST("", categoryHandler.CreateCategory)
		categoryProtectedRoutes.PUT("/:id", categoryHandler.UpdateCategory)
		categoryProtectedRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
	}
}
