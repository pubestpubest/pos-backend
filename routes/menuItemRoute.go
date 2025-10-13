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
	menuItemHandler := menuItemHandler.NewMenuItemHandler(menuItemUsecase)

	// Public routes for customers
	menuItemPublicRoutes := v1.Group("/menu-items")
	{
		menuItemPublicRoutes.GET("", menuItemHandler.GetAllMenuItems)
		menuItemPublicRoutes.GET("/modifiers", menuItemHandler.GetAvailableModifiers)
		menuItemPublicRoutes.GET("/:id", menuItemHandler.GetMenuItemByID)
	}

	// Protected routes for staff/admin
	menuItemProtectedRoutes := v1.Group("/menu-items")
	menuItemProtectedRoutes.Use(middlewares.AuthMiddleware())
	{
		menuItemProtectedRoutes.POST("", menuItemHandler.CreateMenuItem)
		menuItemProtectedRoutes.PUT("/:id", menuItemHandler.UpdateMenuItem)
		menuItemProtectedRoutes.DELETE("/:id", menuItemHandler.DeleteMenuItem)
	}
}
