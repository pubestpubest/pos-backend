package repository

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/pubestpubest/pos-backend/domain"
	"github.com/pubestpubest/pos-backend/models"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.
		Preload("Table").
		Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // Include soft-deleted menu items for historical orders
		}).
		Preload("Items.Modifiers.Modifier").
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "[OrderRepository.GetAllOrders]: Error querying database")
	}
	return orders, nil
}

func (r *orderRepository) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.Where("id = ?", id).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[OrderRepository.GetOrderByID]: Order not found")
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetOrderByID]: Error querying database")
	}
	return &order, nil
}

func (r *orderRepository) GetOrderWithItems(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.db.
		Preload("Table").
		Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // Include soft-deleted menu items for historical orders
		}).
		Preload("Items.Modifiers.Modifier").
		Where("id = ?", id).
		First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[OrderRepository.GetOrderWithItems]: Order not found")
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetOrderWithItems]: Error querying database")
	}
	return &order, nil
}

func (r *orderRepository) GetOrdersByTable(tableID uuid.UUID) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.
		Preload("Table").
		Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // Include soft-deleted menu items for historical orders
		}).
		Preload("Items.Modifiers.Modifier").
		Where("table_id = ?", tableID).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "[OrderRepository.GetOrdersByTable]: Error querying database")
	}
	return orders, nil
}

func (r *orderRepository) GetOrdersByStatus(status string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.
		Preload("Table").
		Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // Include soft-deleted menu items for historical orders
		}).
		Preload("Items.Modifiers.Modifier").
		Where("status = ?", status).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, errors.Wrap(err, "[OrderRepository.GetOrdersByStatus]: Error querying database")
	}
	return orders, nil
}

func (r *orderRepository) CreateOrder(order *models.Order) error {
	if err := r.db.Create(order).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.CreateOrder]: Error creating order")
	}
	return nil
}

func (r *orderRepository) UpdateOrder(order *models.Order) error {
	if err := r.db.Save(order).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.UpdateOrder]: Error updating order")
	}
	return nil
}

func (r *orderRepository) CreateOrderItem(item *models.OrderItem) error {
	if err := r.db.Create(item).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.CreateOrderItem]: Error creating order item")
	}
	return nil
}

func (r *orderRepository) UpdateOrderItem(item *models.OrderItem) error {
	if err := r.db.Save(item).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.UpdateOrderItem]: Error updating order item")
	}
	return nil
}

func (r *orderRepository) DeleteOrderItem(id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&models.OrderItem{}).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.DeleteOrderItem]: Error deleting order item")
	}
	return nil
}

func (r *orderRepository) GetOrderItemByID(id uuid.UUID) (*models.OrderItem, error) {
	var orderItem models.OrderItem
	if err := r.db.
		Preload("MenuItem", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped() // Include soft-deleted menu items for historical orders
		}).
		Preload("Modifiers.Modifier").
		Where("id = ?", id).
		First(&orderItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[OrderRepository.GetOrderItemByID]: Order item not found")
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetOrderItemByID]: Error querying database")
	}
	return &orderItem, nil
}

func (r *orderRepository) GetMenuItemByID(id uuid.UUID) (*models.MenuItem, error) {
	var menuItem models.MenuItem
	if err := r.db.Where("id = ?", id).First(&menuItem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[OrderRepository.GetMenuItemByID]: Menu item not found")
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetMenuItemByID]: Error querying database")
	}
	return &menuItem, nil
}

func (r *orderRepository) GetModifierByID(id uuid.UUID) (*models.Modifier, error) {
	var modifier models.Modifier
	if err := r.db.Where("id = ?", id).First(&modifier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[OrderRepository.GetModifierByID]: Modifier not found")
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetModifierByID]: Error querying database")
	}
	return &modifier, nil
}

func (r *orderRepository) CreateOrderItemModifier(modifier *models.OrderItemModifier) error {
	if err := r.db.Create(modifier).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.CreateOrderItemModifier]: Error creating order item modifier")
	}
	return nil
}

func (r *orderRepository) DeleteOrderItemModifiers(orderItemID uuid.UUID) error {
	if err := r.db.Where("order_item_id = ?", orderItemID).Delete(&models.OrderItemModifier{}).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.DeleteOrderItemModifiers]: Error deleting order item modifiers")
	}
	return nil
}

func (r *orderRepository) GetTableByID(id uuid.UUID) (*models.DiningTable, error) {
	var table models.DiningTable
	if err := r.db.Where("id = ?", id).First(&table).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Wrap(err, "[OrderRepository.GetTableByID]: Table not found")
		}
		return nil, errors.Wrap(err, "[OrderRepository.GetTableByID]: Error querying database")
	}
	return &table, nil
}

func (r *orderRepository) UpdateTable(table *models.DiningTable) error {
	if err := r.db.Save(table).Error; err != nil {
		return errors.Wrap(err, "[OrderRepository.UpdateTable]: Error updating table")
	}
	return nil
}
