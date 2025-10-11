package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pubestpubest/pos-backend/database"
	paymentHandler "github.com/pubestpubest/pos-backend/feature/payment/delivery"
	paymentRepository "github.com/pubestpubest/pos-backend/feature/payment/repository"
	paymentUsecase "github.com/pubestpubest/pos-backend/feature/payment/usecase"
	"github.com/pubestpubest/pos-backend/middlewares"
)

func PaymentRoutes(v1 *gin.RouterGroup) {
	paymentRepository := paymentRepository.NewPaymentRepository(database.DB)
	paymentUsecase := paymentUsecase.NewPaymentUsecase(paymentRepository)
	paymentHandler := paymentHandler.NewPaymentHandler(paymentUsecase)

	paymentRoutes := v1.Group("/payments")
	paymentRoutes.Use(middlewares.AuthMiddleware())
	{
		paymentRoutes.GET("", paymentHandler.GetAllPayments)
		paymentRoutes.GET("/:id", paymentHandler.GetPaymentByID)
		paymentRoutes.POST("", paymentHandler.ProcessPayment)
		paymentRoutes.GET("/methods", paymentHandler.GetPaymentMethods)
	}

	// Order-specific payment routes
	orderPaymentRoutes := v1.Group("/orders/:order_id/payments")
	orderPaymentRoutes.Use(middlewares.AuthMiddleware())
	{
		orderPaymentRoutes.GET("", paymentHandler.GetPaymentsByOrder)
	}
}
