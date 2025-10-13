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

	// Public routes for customers
	paymentPublicRoutes := v1.Group("/payments")
	{
		paymentPublicRoutes.POST("", paymentHandler.ProcessPayment)
		paymentPublicRoutes.GET("/methods", paymentHandler.GetPaymentMethods)
	}

	// Protected routes for staff/admin
	paymentProtectedRoutes := v1.Group("/payments")
	paymentProtectedRoutes.Use(middlewares.AuthMiddleware())
	{
		paymentProtectedRoutes.GET("", paymentHandler.GetAllPayments)
		paymentProtectedRoutes.GET("/:id", paymentHandler.GetPaymentByID)
	}

	// Order-specific payment routes (public for customers to view their order payments)
	orderPaymentPublicRoutes := v1.Group("/orders/:id/payments")
	{
		orderPaymentPublicRoutes.GET("", paymentHandler.GetPaymentsByOrder)
	}
}
