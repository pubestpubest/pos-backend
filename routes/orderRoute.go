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

	// Public routes for customers
	orderPublicRoutes := v1.Group("/orders")
	{
		orderPublicRoutes.GET("/:id", orderHandler.GetOrderByID)
		orderPublicRoutes.POST("", orderHandler.CreateOrder)
		orderPublicRoutes.POST("/:id/items", orderHandler.AddItemToOrder)
		orderPublicRoutes.DELETE("/:id/items/:item_id", orderHandler.RemoveItemFromOrder)
		orderPublicRoutes.PUT("/:id/items/:item_id/quantity", orderHandler.UpdateOrderItemQuantity)
	}

	// Protected routes for staff/admin
	orderProtectedRoutes := v1.Group("/orders")
	orderProtectedRoutes.Use(middlewares.AuthMiddleware())
	{
		orderProtectedRoutes.GET("", orderHandler.GetAllOrders)
		orderProtectedRoutes.GET("/open", orderHandler.GetOpenOrders)
		orderProtectedRoutes.PUT("/:id/close", orderHandler.CloseOrder)
		orderProtectedRoutes.PUT("/:id/void", orderHandler.VoidOrder)
	}

	// Table-specific routes (protected)
	tableOrderRoutes := v1.Group("/tables/:id/orders")
	tableOrderRoutes.Use(middlewares.AuthMiddleware())
	{
		tableOrderRoutes.GET("", orderHandler.GetOrdersByTable)
	}
}
