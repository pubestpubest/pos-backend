package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/utils"
	log "github.com/sirupsen/logrus"
)

type paymentHandler struct {
	paymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(paymentUsecase domain.PaymentUsecase) *paymentHandler {
	return &paymentHandler{paymentUsecase: paymentUsecase}
}

func (h *paymentHandler) GetAllPayments(c *gin.Context) {
	payments, err := h.paymentUsecase.GetAllPayments()
	if err != nil {
		err = errors.Wrap(err, "[PaymentHandler.GetAllPayments]: Error getting payments")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func (h *paymentHandler) GetPaymentByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
		return
	}

	payment, err := h.paymentUsecase.GetPaymentByID(id)
	if err != nil {
		err = errors.Wrap(err, "[PaymentHandler.GetPaymentByID]: Error getting payment")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func (h *paymentHandler) GetPaymentsByOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	payments, err := h.paymentUsecase.GetPaymentsByOrder(orderID)
	if err != nil {
		err = errors.Wrap(err, "[PaymentHandler.GetPaymentsByOrder]: Error getting payments")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func (h *paymentHandler) ProcessPayment(c *gin.Context) {
	var req request.PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	payment, err := h.paymentUsecase.ProcessPayment(&req)
	if err != nil {
		err = errors.Wrap(err, "[PaymentHandler.ProcessPayment]: Error processing payment")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusCreated, payment)
}

func (h *paymentHandler) GetPaymentMethods(c *gin.Context) {
	methods, err := h.paymentUsecase.GetPaymentMethods()
	if err != nil {
		err = errors.Wrap(err, "[PaymentHandler.GetPaymentMethods]: Error getting payment methods")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, methods)
}
