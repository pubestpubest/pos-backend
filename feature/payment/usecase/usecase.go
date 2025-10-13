package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/constant"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
	"github.com/pubestpubest/pos-backend/utils"
)

type paymentUsecase struct {
	paymentRepository domain.PaymentRepository
}

func NewPaymentUsecase(paymentRepository domain.PaymentRepository) domain.PaymentUsecase {
	return &paymentUsecase{paymentRepository: paymentRepository}
}

func (u *paymentUsecase) GetAllPayments() ([]*response.PaymentResponse, error) {
	payments, err := u.paymentRepository.GetAllPayments()
	if err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.GetAllPayments]: Error getting payments")
	}

	return u.buildPaymentResponses(payments), nil
}

func (u *paymentUsecase) GetPaymentByID(id uuid.UUID) (*response.PaymentResponse, error) {
	payment, err := u.paymentRepository.GetPaymentByID(id)
	if err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.GetPaymentByID]: Error getting payment")
	}

	return u.buildPaymentResponse(payment), nil
}

func (u *paymentUsecase) GetPaymentsByOrder(orderID uuid.UUID) ([]*response.PaymentResponse, error) {
	payments, err := u.paymentRepository.GetPaymentsByOrder(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.GetPaymentsByOrder]: Error getting payments")
	}

	return u.buildPaymentResponses(payments), nil
}

func (u *paymentUsecase) ProcessPayment(req *request.PaymentRequest) (*response.PaymentResponse, error) {
	// Validate order exists
	order, err := u.paymentRepository.GetOrderByID(req.OrderID)
	if err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Order not found")
	}

	// Check if order is open
	if *order.Status != constant.OrderStatusOpen {
		return nil, errors.New("[PaymentUsecase.ProcessPayment]: Can only pay for open orders")
	}

	// Get total already paid
	totalPaid, err := u.paymentRepository.GetTotalPaidForOrder(req.OrderID)
	if err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Error checking payment status")
	}

	// Validate payment amount
	orderTotal := utils.DerefInt64(order.TotalBaht)
	if totalPaid+req.AmountBaht > orderTotal {
		return nil, errors.New("[PaymentUsecase.ProcessPayment]: Payment amount exceeds order total")
	}

	// Create payment
	payment := &models.Payment{
		OrderID:     req.OrderID,
		Method:      &req.Method,
		AmountBaht:  req.AmountBaht,
		Currency:    utils.Ptr(constant.PaymentCurrencyTHB),
		Provider:    req.Provider,
		ProviderRef: req.ProviderRef,
		Status:      utils.Ptr(constant.PaymentStatusSucceeded),
	}

	if err := u.paymentRepository.CreatePayment(payment); err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Error processing payment")
	}

	// Calculate new total paid
	newTotalPaid := totalPaid + req.AmountBaht

	// If order is fully paid, close the order and free the table
	if newTotalPaid >= orderTotal {
		// Close the order
		order.Status = utils.Ptr(constant.OrderStatusPaid)
		order.ClosedAt = utils.Ptr(time.Now())
		if err := u.paymentRepository.UpdateOrder(order); err != nil {
			return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Error closing order")
		}

		// Free the table if order has a table
		if order.TableID != nil {
			table, err := u.paymentRepository.GetTableByID(*order.TableID)
			if err != nil {
				return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Error getting table")
			}

			table.Status = utils.Ptr(constant.TableStatusFree)
			if err := u.paymentRepository.UpdateTable(table); err != nil {
				return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Error freeing table")
			}
		}
	}

	// Reload payment with order
	paymentWithOrder, err := u.paymentRepository.GetPaymentByID(payment.ID)
	if err != nil {
		return nil, errors.Wrap(err, "[PaymentUsecase.ProcessPayment]: Error retrieving payment")
	}

	return u.buildPaymentResponse(paymentWithOrder), nil
}

func (u *paymentUsecase) GetPaymentMethods() ([]*response.PaymentMethodResponse, error) {
	// Return static list of payment methods
	methods := []*response.PaymentMethodResponse{
		{
			Code: constant.PaymentMethodCash,
			Name: "Cash",
		},
		{
			Code: constant.PaymentMethodCard,
			Name: "Credit/Debit Card",
		},
		{
			Code: constant.PaymentMethodPromptpay,
			Name: "PromptPay",
		},
	}

	return methods, nil
}

// Helper function to build payment response
func (u *paymentUsecase) buildPaymentResponse(payment *models.Payment) *response.PaymentResponse {
	return &response.PaymentResponse{
		ID:          payment.ID,
		OrderID:     payment.OrderID,
		Method:      utils.DerefString(payment.Method),
		AmountBaht:  payment.AmountBaht,
		Currency:    utils.DerefString(payment.Currency),
		Provider:    utils.DerefString(payment.Provider),
		ProviderRef: utils.DerefString(payment.ProviderRef),
		Status:      utils.DerefString(payment.Status),
		CreatedAt:   payment.CreatedAt,
	}
}

// Helper function to build multiple payment responses
func (u *paymentUsecase) buildPaymentResponses(payments []*models.Payment) []*response.PaymentResponse {
	responses := make([]*response.PaymentResponse, len(payments))
	for i, payment := range payments {
		responses[i] = u.buildPaymentResponse(payment)
	}
	return responses
}
