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

type orderUsecase struct {
	orderRepository domain.OrderRepository
}

func NewOrderUsecase(orderRepository domain.OrderRepository) domain.OrderUsecase {
	return &orderUsecase{orderRepository: orderRepository}
}

func (u *orderUsecase) GetAllOrders() ([]*response.OrderResponse, error) {
	orders, err := u.orderRepository.GetAllOrders()
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.GetAllOrders]: Error getting orders")
	}

	return u.buildOrderResponses(orders), nil
}

func (u *orderUsecase) GetOrderByID(id uuid.UUID) (*response.OrderResponse, error) {
	order, err := u.orderRepository.GetOrderWithItems(id)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.GetOrderByID]: Error getting order")
	}

	return u.buildOrderResponse(order), nil
}

func (u *orderUsecase) GetOrdersByTable(tableID uuid.UUID) ([]*response.OrderResponse, error) {
	orders, err := u.orderRepository.GetOrdersByTable(tableID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.GetOrdersByTable]: Error getting orders")
	}

	return u.buildOrderResponses(orders), nil
}

func (u *orderUsecase) GetOpenOrders() ([]*response.OrderResponse, error) {
	orders, err := u.orderRepository.GetOrdersByStatus(constant.OrderStatusOpen)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.GetOpenOrders]: Error getting open orders")
	}

	return u.buildOrderResponses(orders), nil
}

func (u *orderUsecase) CreateOrder(req *request.OrderCreateRequest) (*response.OrderResponse, error) {
	// Validate table exists
	_, err := u.orderRepository.GetTableByID(req.TableID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.CreateOrder]: Invalid table ID")
	}

	order := &models.Order{
		TableID:      &req.TableID,
		OpenedBy:     req.OpenedBy,
		Source:       &req.Source,
		Status:       utils.Ptr(constant.OrderStatusOpen),
		SubtotalBaht: utils.PtrI64(0),
		DiscountBaht: utils.PtrI64(0),
		TotalBaht:    utils.PtrI64(0),
		Note:         req.Note,
	}

	if err := u.orderRepository.CreateOrder(order); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.CreateOrder]: Error creating order")
	}

	// Get order with items
	orderWithItems, err := u.orderRepository.GetOrderWithItems(order.ID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.CreateOrder]: Error retrieving created order")
	}

	return u.buildOrderResponse(orderWithItems), nil
}

func (u *orderUsecase) AddItemToOrder(orderID uuid.UUID, req *request.AddOrderItemRequest) (*response.OrderResponse, error) {
	// Get order
	order, err := u.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Order not found")
	}

	// Check if order is open
	if *order.Status != constant.OrderStatusOpen {
		return nil, errors.New("[OrderUsecase.AddItemToOrder]: Cannot add items to closed order")
	}

	// Get menu item
	menuItem, err := u.orderRepository.GetMenuItemByID(req.MenuItemID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Menu item not found")
	}

	// Calculate unit price (base price)
	unitPrice := utils.DerefInt64(menuItem.PriceBaht)

	// Calculate modifier total
	modifierTotal := int64(0)
	for _, modifierID := range req.ModifierIDs {
		modifier, err := u.orderRepository.GetModifierByID(modifierID)
		if err != nil {
			return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Modifier not found")
		}
		modifierTotal += utils.DerefInt64(modifier.PriceDeltaBaht)
	}

	// Calculate line total = (base price + modifier total) * quantity
	lineTotal := (unitPrice + modifierTotal) * int64(req.Quantity)

	// Create order item
	orderItem := &models.OrderItem{
		OrderID:       orderID,
		MenuItemID:    req.MenuItemID,
		Quantity:      req.Quantity,
		UnitPriceBaht: unitPrice,
		LineTotalBaht: lineTotal,
		Note:          req.Note,
	}

	if err := u.orderRepository.CreateOrderItem(orderItem); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Error creating order item")
	}

	// Add modifiers
	for _, modifierID := range req.ModifierIDs {
		modifier, _ := u.orderRepository.GetModifierByID(modifierID)
		orderItemModifier := &models.OrderItemModifier{
			OrderItemID:    orderItem.ID,
			ModifierID:     modifierID,
			PriceDeltaBaht: modifier.PriceDeltaBaht,
		}
		if err := u.orderRepository.CreateOrderItemModifier(orderItemModifier); err != nil {
			return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Error adding modifier")
		}
	}

	// Recalculate order total
	if err := u.recalculateOrderTotal(orderID); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Error recalculating total")
	}

	// Get updated order
	updatedOrder, err := u.orderRepository.GetOrderWithItems(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.AddItemToOrder]: Error retrieving updated order")
	}

	return u.buildOrderResponse(updatedOrder), nil
}

func (u *orderUsecase) RemoveItemFromOrder(orderID uuid.UUID, itemID uuid.UUID) (*response.OrderResponse, error) {
	// Get order
	order, err := u.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.RemoveItemFromOrder]: Order not found")
	}

	// Check if order is open
	if *order.Status != constant.OrderStatusOpen {
		return nil, errors.New("[OrderUsecase.RemoveItemFromOrder]: Cannot remove items from closed order")
	}

	// Get order item
	orderItem, err := u.orderRepository.GetOrderItemByID(itemID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.RemoveItemFromOrder]: Order item not found")
	}

	// Verify item belongs to order
	if orderItem.OrderID != orderID {
		return nil, errors.New("[OrderUsecase.RemoveItemFromOrder]: Order item does not belong to this order")
	}

	// Delete modifiers first
	if err := u.orderRepository.DeleteOrderItemModifiers(itemID); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.RemoveItemFromOrder]: Error deleting modifiers")
	}

	// Delete order item
	if err := u.orderRepository.DeleteOrderItem(itemID); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.RemoveItemFromOrder]: Error deleting order item")
	}

	// Recalculate order total
	if err := u.recalculateOrderTotal(orderID); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.RemoveItemFromOrder]: Error recalculating total")
	}

	// Get updated order
	updatedOrder, err := u.orderRepository.GetOrderWithItems(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.RemoveItemFromOrder]: Error retrieving updated order")
	}

	return u.buildOrderResponse(updatedOrder), nil
}

func (u *orderUsecase) UpdateOrderItemQuantity(orderID uuid.UUID, itemID uuid.UUID, quantity int) (*response.OrderResponse, error) {
	// Get order
	order, err := u.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.UpdateOrderItemQuantity]: Order not found")
	}

	// Check if order is open
	if *order.Status != constant.OrderStatusOpen {
		return nil, errors.New("[OrderUsecase.UpdateOrderItemQuantity]: Cannot update items in closed order")
	}

	// Get order item with modifiers
	orderItem, err := u.orderRepository.GetOrderItemByID(itemID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.UpdateOrderItemQuantity]: Order item not found")
	}

	// Verify item belongs to order
	if orderItem.OrderID != orderID {
		return nil, errors.New("[OrderUsecase.UpdateOrderItemQuantity]: Order item does not belong to this order")
	}

	// Calculate modifier total
	modifierTotal := int64(0)
	for _, mod := range orderItem.Modifiers {
		modifierTotal += utils.DerefInt64(mod.PriceDeltaBaht)
	}

	// Update quantity and line total
	orderItem.Quantity = quantity
	orderItem.LineTotalBaht = (orderItem.UnitPriceBaht + modifierTotal) * int64(quantity)

	if err := u.orderRepository.UpdateOrderItem(orderItem); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.UpdateOrderItemQuantity]: Error updating order item")
	}

	// Recalculate order total
	if err := u.recalculateOrderTotal(orderID); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.UpdateOrderItemQuantity]: Error recalculating total")
	}

	// Get updated order
	updatedOrder, err := u.orderRepository.GetOrderWithItems(orderID)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.UpdateOrderItemQuantity]: Error retrieving updated order")
	}

	return u.buildOrderResponse(updatedOrder), nil
}

func (u *orderUsecase) CloseOrder(id uuid.UUID) (*response.OrderResponse, error) {
	order, err := u.orderRepository.GetOrderWithItems(id)
	if err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.CloseOrder]: Order not found")
	}

	// Check if order is already closed
	if *order.Status != constant.OrderStatusOpen {
		return nil, errors.New("[OrderUsecase.CloseOrder]: Order is already closed")
	}

	// Update order status
	now := time.Now()
	order.Status = utils.Ptr(constant.OrderStatusPaid)
	order.ClosedAt = &now

	if err := u.orderRepository.UpdateOrder(order); err != nil {
		return nil, errors.Wrap(err, "[OrderUsecase.CloseOrder]: Error closing order")
	}

	return u.buildOrderResponse(order), nil
}

func (u *orderUsecase) VoidOrder(id uuid.UUID) error {
	order, err := u.orderRepository.GetOrderByID(id)
	if err != nil {
		return errors.Wrap(err, "[OrderUsecase.VoidOrder]: Order not found")
	}

	// Update order status
	order.Status = utils.Ptr(constant.OrderStatusVoid)

	if err := u.orderRepository.UpdateOrder(order); err != nil {
		return errors.Wrap(err, "[OrderUsecase.VoidOrder]: Error voiding order")
	}

	return nil
}

// Helper function to recalculate order total
func (u *orderUsecase) recalculateOrderTotal(orderID uuid.UUID) error {
	order, err := u.orderRepository.GetOrderWithItems(orderID)
	if err != nil {
		return err
	}

	subtotal := int64(0)
	for _, item := range order.Items {
		subtotal += item.LineTotalBaht
	}

	discount := utils.DerefInt64(order.DiscountBaht)
	total := subtotal - discount

	order.SubtotalBaht = &subtotal
	order.TotalBaht = &total

	return u.orderRepository.UpdateOrder(order)
}

// Helper function to build order response
func (u *orderUsecase) buildOrderResponse(order *models.Order) *response.OrderResponse {
	items := make([]response.OrderItemResponse, len(order.Items))
	for i, item := range order.Items {
		modifiers := make([]response.OrderItemModifierResponse, len(item.Modifiers))
		for j, mod := range item.Modifiers {
			modifiers[j] = response.OrderItemModifierResponse{
				ModifierID:     mod.ModifierID,
				ModifierName:   utils.DerefString(mod.Modifier.Name),
				PriceDeltaBaht: utils.DerefInt64(mod.PriceDeltaBaht),
			}
		}

		items[i] = response.OrderItemResponse{
			ID:            item.ID,
			MenuItemID:    item.MenuItemID,
			MenuItemName:  utils.DerefString(item.MenuItem.Name),
			Quantity:      item.Quantity,
			UnitPriceBaht: item.UnitPriceBaht,
			LineTotalBaht: item.LineTotalBaht,
			Note:          utils.DerefString(item.Note),
			Modifiers:     modifiers,
		}
	}

	return &response.OrderResponse{
		ID:           order.ID,
		TableID:      utils.DerefUUID(order.TableID),
		TableName:    utils.DerefString(order.Table.Name),
		OpenedBy:     order.OpenedBy,
		Source:       utils.DerefString(order.Source),
		Status:       utils.DerefString(order.Status),
		SubtotalBaht: utils.DerefInt64(order.SubtotalBaht),
		DiscountBaht: utils.DerefInt64(order.DiscountBaht),
		TotalBaht:    utils.DerefInt64(order.TotalBaht),
		Note:         utils.DerefString(order.Note),
		CreatedAt:    order.CreatedAt,
		ClosedAt:     order.ClosedAt,
		Items:        items,
	}
}

// Helper function to build multiple order responses
func (u *orderUsecase) buildOrderResponses(orders []*models.Order) []*response.OrderResponse {
	responses := make([]*response.OrderResponse, len(orders))
	for i, order := range orders {
		responses[i] = u.buildOrderResponse(order)
	}
	return responses
}
