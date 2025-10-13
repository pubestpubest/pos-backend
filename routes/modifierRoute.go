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

	// Public routes for customers
	modifierPublicRoutes := v1.Group("/modifiers")
	{
		modifierPublicRoutes.GET("", modifierHandler.GetAllModifiers)
		modifierPublicRoutes.GET("/:id", modifierHandler.GetModifierByID)
	}

	// Protected routes for staff/admin
	modifierProtectedRoutes := v1.Group("/modifiers")
	modifierProtectedRoutes.Use(middlewares.AuthMiddleware())
	{
		modifierProtectedRoutes.POST("", modifierHandler.CreateModifier)
		modifierProtectedRoutes.PUT("/:id", modifierHandler.UpdateModifier)
		modifierProtectedRoutes.DELETE("/:id", modifierHandler.DeleteModifier)
	}
}
