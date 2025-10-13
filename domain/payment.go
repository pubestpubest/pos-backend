package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// Payment domain - manages order payments
type PaymentUsecase interface {
	GetAllPayments() ([]*response.PaymentResponse, error)
	GetPaymentByID(id uuid.UUID) (*response.PaymentResponse, error)
	GetPaymentsByOrder(orderID uuid.UUID) ([]*response.PaymentResponse, error)
	ProcessPayment(req *request.PaymentRequest) (*response.PaymentResponse, error)
	GetPaymentMethods() ([]*response.PaymentMethodResponse, error)
}

type PaymentRepository interface {
	GetAllPayments() ([]*models.Payment, error)
	GetPaymentByID(id uuid.UUID) (*models.Payment, error)
	GetPaymentsByOrder(orderID uuid.UUID) ([]*models.Payment, error)
	CreatePayment(payment *models.Payment) error
	UpdatePayment(payment *models.Payment) error
	GetTotalPaidForOrder(orderID uuid.UUID) (int64, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	GetTableByID(id uuid.UUID) (*models.DiningTable, error)
	UpdateTable(table *models.DiningTable) error
}
