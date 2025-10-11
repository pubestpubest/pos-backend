package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	roleHandler "github.com/pubestpubest/pos-backend/feature/role/delivery"
	roleRepository "github.com/pubestpubest/pos-backend/feature/role/repository"
	roleUsecase "github.com/pubestpubest/pos-backend/feature/role/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func RoleRoutes(v1 *gin.RouterGroup) {
	roleRepository := roleRepository.NewRoleRepository(database.DB)
	roleUsecase := roleUsecase.NewRoleUsecase(roleRepository)
	roleHandler := roleHandler.NewRoleHandler(roleUsecase)

	roleRoutes := v1.Group("/roles")
	roleRoutes.Use(middlewares.AuthMiddleware())
	{
		roleRoutes.GET("", roleHandler.GetAllRoles)
		roleRoutes.GET("/:id", roleHandler.GetRoleWithPermissions)
	}
}
