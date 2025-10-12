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

type orderHandler struct {
	orderUsecase domain.OrderUsecase
}

func NewOrderHandler(orderUsecase domain.OrderUsecase) *orderHandler {
	return &orderHandler{orderUsecase: orderUsecase}
}

func (h *orderHandler) GetAllOrders(c *gin.Context) {
	orders, err := h.orderUsecase.GetAllOrders()
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.GetAllOrders]: Error getting orders")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *orderHandler) GetOrderByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderUsecase.GetOrderByID(id)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.GetOrderByID]: Error getting order")
		log.Warn(err)
		c.JSON(http.StatusNotFound, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *orderHandler) GetOrdersByTable(c *gin.Context) {
	tableID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
		return
	}

	orders, err := h.orderUsecase.GetOrdersByTable(tableID)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.GetOrdersByTable]: Error getting orders")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *orderHandler) GetOpenOrders(c *gin.Context) {
	orders, err := h.orderUsecase.GetOpenOrders()
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.GetOpenOrders]: Error getting open orders")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *orderHandler) CreateOrder(c *gin.Context) {
	var req request.OrderCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	order, err := h.orderUsecase.CreateOrder(&req)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.CreateOrder]: Error creating order")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *orderHandler) AddItemToOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	var req request.AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	order, err := h.orderUsecase.AddItemToOrder(orderID, &req)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.AddItemToOrder]: Error adding item to order")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *orderHandler) RemoveItemFromOrder(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	order, err := h.orderUsecase.RemoveItemFromOrder(orderID, itemID)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.RemoveItemFromOrder]: Error removing item from order")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *orderHandler) UpdateOrderItemQuantity(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	itemID, err := uuid.Parse(c.Param("item_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var req request.UpdateOrderItemQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	order, err := h.orderUsecase.UpdateOrderItemQuantity(orderID, itemID, req.Quantity)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.UpdateOrderItemQuantity]: Error updating item quantity")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *orderHandler) CloseOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	order, err := h.orderUsecase.CloseOrder(id)
	if err != nil {
		err = errors.Wrap(err, "[OrderHandler.CloseOrder]: Error closing order")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *orderHandler) VoidOrder(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}

	if err := h.orderUsecase.VoidOrder(id); err != nil {
		err = errors.Wrap(err, "[OrderHandler.VoidOrder]: Error voiding order")
		log.Warn(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.StandardError(err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order voided successfully"})
}
