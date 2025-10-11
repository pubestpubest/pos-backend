package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	modifierHandler "github.com/pubestpubest/pos-backend/feature/modifier/delivery"
	modifierRepository "github.com/pubestpubest/pos-backend/feature/modifier/repository"
	modifierUsecase "github.com/pubestpubest/pos-backend/feature/modifier/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func ModifierRoutes(v1 *gin.RouterGroup) {
	modifierRepository := modifierRepository.NewModifierRepository(database.DB)
	modifierUsecase := modifierUsecase.NewModifierUsecase(modifierRepository)
	modifierHandler := modifierHandler.NewModifierHandler(modifierUsecase)

	modifierRoutes := v1.Group("/modifiers")
	modifierRoutes.Use(middlewares.AuthMiddleware())
	{
		modifierRoutes.GET("", modifierHandler.GetAllModifiers)
		modifierRoutes.GET("/:id", modifierHandler.GetModifierByID)
		modifierRoutes.POST("", modifierHandler.CreateModifier)
		modifierRoutes.PUT("/:id", modifierHandler.UpdateModifier)
		modifierRoutes.DELETE("/:id", modifierHandler.DeleteModifier)
	}
}
