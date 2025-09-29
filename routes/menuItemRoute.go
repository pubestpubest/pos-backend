package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	menuItemHandler "github.com/pubestpubest/pos-backend/feature/menuItem/delivery"
	menuItemRepository "github.com/pubestpubest/pos-backend/feature/menuItem/repository"
	menuItemUsecase "github.com/pubestpubest/pos-backend/feature/menuItem/usecase"
)

func MenuItemRoutes(v1 *gin.RouterGroup) {

	menuItemRepository := menuItemRepository.NewMenuItemRepository(database.DB)
	menuItemUsecase := menuItemUsecase.NewMenuItemUsecase(menuItemRepository)
	menuItemHandler := menuItemHandler.NewMenuItemHandler(menuItemUsecase)

	menuItemRoutes := v1.Group("/menu-items")

	menuItemRoutes.GET("/", menuItemHandler.GetAllMenuItems)
}
