package domain

import (
	"github.com/google/uuid"
	"github.com/pubestpubest/pos-backend/models"
	"github.com/pubestpubest/pos-backend/request"
	"github.com/pubestpubest/pos-backend/response"
)

// Order domain - manages customer orders and order items
type OrderUsecase interface {
	GetAllOrders() ([]*response.OrderResponse, error)
	GetOrderByID(id uuid.UUID) (*response.OrderResponse, error)
	GetOrdersByTable(tableID uuid.UUID) ([]*response.OrderResponse, error)
	GetOpenOrders() ([]*response.OrderResponse, error)
	CreateOrder(req *request.OrderCreateRequest) (*response.OrderResponse, error)
	AddItemToOrder(orderID uuid.UUID, req *request.AddOrderItemRequest) (*response.OrderResponse, error)
	RemoveItemFromOrder(orderID uuid.UUID, itemID uuid.UUID) (*response.OrderResponse, error)
	UpdateOrderItemQuantity(orderID uuid.UUID, itemID uuid.UUID, quantity int) (*response.OrderResponse, error)
	CloseOrder(id uuid.UUID) (*response.OrderResponse, error)
	VoidOrder(id uuid.UUID) error
}

type OrderRepository interface {
	GetAllOrders() ([]*models.Order, error)
	GetOrderByID(id uuid.UUID) (*models.Order, error)
	GetOrderWithItems(id uuid.UUID) (*models.Order, error)
	GetOrdersByTable(tableID uuid.UUID) ([]*models.Order, error)
	GetOrdersByStatus(status string) ([]*models.Order, error)
	CreateOrder(order *models.Order) error
	UpdateOrder(order *models.Order) error
	CreateOrderItem(item *models.OrderItem) error
	UpdateOrderItem(item *models.OrderItem) error
	DeleteOrderItem(id uuid.UUID) error
	GetOrderItemByID(id uuid.UUID) (*models.OrderItem, error)
	GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error)
	GetModifierByID(id uuid.UUID) (*models.Modifier, error)
	CreateOrderItemModifier(modifier *models.OrderItemModifier) error
	DeleteOrderItemModifiers(orderItemID uuid.UUID) error
	GetTableByID(id uuid.UUID) (*models.DiningTable, error)
	UpdateTable(table *models.DiningTable) error
}
