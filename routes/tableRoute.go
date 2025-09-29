package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	tableHandler "github.com/pubestpubest/pos-backend/feature/table/delivery"
	tableRepository "github.com/pubestpubest/pos-backend/feature/table/repository"
	tableUsecase "github.com/pubestpubest/pos-backend/feature/table/usecase"
)

func TableRoutes(v1 *gin.RouterGroup) {

	tableRepository := tableRepository.NewTableRepository(database.DB)
	tableUsecase := tableUsecase.NewTableUsecase(tableRepository)
	tableHandler := tableHandler.NewTableHandler(tableUsecase)

	tableRoutes := v1.Group("/tables")

	tableRoutes.GET("/", tableHandler.GetAllTables)
}
