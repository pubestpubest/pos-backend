package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	orderHandler "github.com/pubestpubest/pos-backend/feature/order/delivery"
	orderRepository "github.com/pubestpubest/pos-backend/feature/order/repository"
	orderUsecase "github.com/pubestpubest/pos-backend/feature/order/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func OrderRoutes(v1 *gin.RouterGroup) {
	orderRepository := orderRepository.NewOrderRepository(database.DB)
	orderUsecase := orderUsecase.NewOrderUsecase(orderRepository)
	orderHandler := orderHandler.NewOrderHandler(orderUsecase)

	orderRoutes := v1.Group("/orders")
	orderRoutes.Use(middlewares.AuthMiddleware())
	{
		orderRoutes.GET("", orderHandler.GetAllOrders)
		orderRoutes.GET("/open", orderHandler.GetOpenOrders)
		orderRoutes.GET("/:id", orderHandler.GetOrderByID)
		orderRoutes.POST("", orderHandler.CreateOrder)
		orderRoutes.POST("/:id/items", orderHandler.AddItemToOrder)
		orderRoutes.DELETE("/:id/items/:item_id", orderHandler.RemoveItemFromOrder)
		orderRoutes.PUT("/:id/items/:item_id/quantity", orderHandler.UpdateOrderItemQuantity)
		orderRoutes.PUT("/:id/close", orderHandler.CloseOrder)
		orderRoutes.PUT("/:id/void", orderHandler.VoidOrder)
	}

	// Table-specific routes
	tableOrderRoutes := v1.Group("/tables/:id/orders")
	tableOrderRoutes.Use(middlewares.AuthMiddleware())
	{
		tableOrderRoutes.GET("", orderHandler.GetOrdersByTable)
	}
}
