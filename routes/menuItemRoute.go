package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	menuItemHandler "github.com/pubestpubest/pos-backend/feature/menuItem/delivery"
	menuItemRepository "github.com/pubestpubest/pos-backend/feature/menuItem/repository"
	menuItemUsecase "github.com/pubestpubest/pos-backend/feature/menuItem/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func MenuItemRoutes(v1 *gin.RouterGroup) {
	menuItemRepository := menuItemRepository.NewMenuItemRepository(database.DB)
	menuItemUsecase := menuItemUsecase.NewMenuItemUsecase(menuItemRepository)
	menuItemHandler := menuItemHandler.NewMenuItemHandler(menuItemUsecase, database.MinioClient)

	// Public routes for customers
	menuItemPublicRoutes := v1.Group("/menu-items")
	{
		menuItemPublicRoutes.GET("", menuItemHandler.GetAllMenuItems)
		menuItemPublicRoutes.GET("/modifiers", menuItemHandler.GetAvailableModifiers)
		menuItemPublicRoutes.GET("/:id", menuItemHandler.GetMenuItemByID)
	}

	// Protected routes for staff/admin (CRUD operations)
	menuItemProtectedRoutes := v1.Group("/menu-items")
	menuItemProtectedRoutes.Use(middlewares.AuthMiddleware())
	{
		menuItemProtectedRoutes.POST("", menuItemHandler.CreateMenuItem)
		menuItemProtectedRoutes.PUT("/:id", menuItemHandler.UpdateMenuItem)
		menuItemProtectedRoutes.DELETE("/:id", menuItemHandler.DeleteMenuItem)
	}

	// Analytics routes (protected - requires authentication)
	analyticsRoutes := v1.Group("/menu-items")
	analyticsRoutes.Use(middlewares.AuthMiddleware())
	{
		analyticsRoutes.GET("/statistics", menuItemHandler.GetAllMenuItemsStatistics)
		analyticsRoutes.GET("/:id/statistics", menuItemHandler.GetMenuItemStatistics)
	}

	// Reports routes (protected - requires authentication)
	reportRoutes := v1.Group("/reports")
	reportRoutes.Use(middlewares.AuthMiddleware())
	{
		reportRoutes.GET("/top-selling-items", menuItemHandler.GetTopSellingItems)
		reportRoutes.GET("/low-selling-items", menuItemHandler.GetLowSellingItems)
	}
}
