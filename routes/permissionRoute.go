package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	permissionHandler "github.com/pubestpubest/pos-backend/feature/permission/delivery"
	permissionRepository "github.com/pubestpubest/pos-backend/feature/permission/repository"
	permissionUsecase "github.com/pubestpubest/pos-backend/feature/permission/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func PermissionRoutes(v1 *gin.RouterGroup) {
	permissionRepository := permissionRepository.NewPermissionRepository(database.DB)
	permissionUsecase := permissionUsecase.NewPermissionUsecase(permissionRepository)
	permissionHandler := permissionHandler.NewPermissionHandler(permissionUsecase)

	permissionRoutes := v1.Group("/permissions")
	permissionRoutes.Use(middlewares.AuthMiddleware())
	{
		permissionRoutes.GET("", permissionHandler.GetAllPermissions)
	}
}
