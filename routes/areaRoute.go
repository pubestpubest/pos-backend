package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	areaHandler "github.com/pubestpubest/pos-backend/feature/area/delivery"
	areaRepository "github.com/pubestpubest/pos-backend/feature/area/repository"
	areaUsecase "github.com/pubestpubest/pos-backend/feature/area/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func AreaRoutes(v1 *gin.RouterGroup) {
	areaRepository := areaRepository.NewAreaRepository(database.DB)
	areaUsecase := areaUsecase.NewAreaUsecase(areaRepository)
	areaHandler := areaHandler.NewAreaHandler(areaUsecase)

	areaRoutes := v1.Group("/areas")
	areaRoutes.Use(middlewares.AuthMiddleware())
	{
		areaRoutes.GET("", areaHandler.GetAllAreas)
		areaRoutes.GET("/:id", areaHandler.GetAreaByID)
		areaRoutes.POST("", areaHandler.CreateArea)
		areaRoutes.PUT("/:id", areaHandler.UpdateArea)
		areaRoutes.DELETE("/:id", areaHandler.DeleteArea)
	}
}
