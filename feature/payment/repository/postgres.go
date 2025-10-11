package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) GetAllPayments() ([]*models.Payment, error) {
	var payments []*models.Payment
	if err := r.db.Preload("Order").Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, errors.Wrap(err, "[PaymentRepository.GetAllPayments]: Error querying database")
	}
	return payments, nil
}

func (r *paymentRepository) GetPaymentByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	if err := r.db.Preload("Order").Where("id = ?", id).First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[PaymentRepository.GetPaymentByID]: Payment not found")
		}
		return nil, errors.Wrap(err, "[PaymentRepository.GetPaymentByID]: Error querying database")
	}
	return &payment, nil
}

func (r *paymentRepository) GetPaymentsByOrder(orderID uuid.UUID) ([]*models.Payment, error) {
	var payments []*models.Payment
	if err := r.db.Where("order_id = ?", orderID).Order("created_at DESC").Find(&payments).Error; err != nil {
		return nil, errors.Wrap(err, "[PaymentRepository.GetPaymentsByOrder]: Error querying database")
	}
	return payments, nil
}

func (r *paymentRepository) CreatePayment(payment *models.Payment) error {
	if err := r.db.Create(payment).Error; err != nil {
		return errors.Wrap(err, "[PaymentRepository.CreatePayment]: Error creating payment")
	}
	return nil
}

func (r *paymentRepository) UpdatePayment(payment *models.Payment) error {
	if err := r.db.Save(payment).Error; err != nil {
		return errors.Wrap(err, "[PaymentRepository.UpdatePayment]: Error updating payment")
	}
	return nil
}

func (r *paymentRepository) GetTotalPaidForOrder(orderID uuid.UUID) (int64, error) {
	var total int64
	if err := r.db.Model(&models.Payment{}).
		Where("order_id = ? AND status = ?", orderID, "succeeded").
		Select("COALESCE(SUM(amount_baht), 0)").
		Scan(&total).Error; err != nil {
		return 0, errors.Wrap(err, "[PaymentRepository.GetTotalPaidForOrder]: Error calculating total")
	}
	return total, nil
}

func (r *paymentRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[PaymentRepository.GetOrderByID]: Order not found")
		}
		return nil, errors.Wrap(err, "[PaymentRepository.GetOrderByID]: Error querying database")
	}
	return &order, nil
}
