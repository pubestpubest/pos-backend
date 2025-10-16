# Order Repository Fix for Soft Delete

## Problem

After implementing soft delete for menu items, orders broke with a nil pointer error:

```
runtime error: invalid memory address or nil pointer dereference
/Users/pubest/gits/pos-backend/feature/order/usecase/usecase.go:346
MenuItemName:  utils.DerefString(item.MenuItem.Name)
```

## Root Cause

When GORM preloads menu items for orders, it automatically applies the soft delete filter. This means:

- ❌ Deleted menu items are excluded from preloads
- ❌ `item.MenuItem` becomes `nil` for items referencing deleted menu items
- ❌ Accessing `item.MenuItem.Name` causes nil pointer panic

**Example:**

```go
// Before fix - causes nil pointer if menu item is deleted
db.Preload("Items.MenuItem")
```

## Solution

Use `Unscoped()` when preloading MenuItem to include soft-deleted items:

```go
// After fix - includes deleted menu items
db.Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
    return db.Unscoped() // Include soft-deleted menu items for historical orders
})
```

## Why This Is Correct

### Historical Data Integrity

Orders are **historical records** that must show what was actually ordered, even if menu items have been deleted since then:

✅ Customer receipts must be accurate  
✅ Financial records must be preserved  
✅ Order history must remain complete  
✅ Audit trails require all data

**Analogy:** If you delete a product from your inventory today, yesterday's invoices should still show that product.

## Changes Made

Updated all order repository methods to use `Unscoped()` for MenuItem preloads:

### 1. GetAllOrders

```go
func (r *orderRepository) GetAllOrders() ([]*models.Order, error) {
    var orders []*models.Order
    if err := r.db.
        Preload("Table").
        Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
            return db.Unscoped() // Include soft-deleted menu items
        }).
        Preload("Items.Modifiers.Modifier").
        Order("created_at DESC").
        Find(&orders).Error; err != nil {
        // ...
    }
    return orders, nil
}
```

### 2. GetOrderWithItems

```go
func (r *orderRepository) GetOrderWithItems(id uuid.UUID) (*models.Order, error) {
    var order models.Order
    if err := r.db.
        Preload("Table").
        Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
            return db.Unscoped() // Include soft-deleted menu items
        }).
        Preload("Items.Modifiers.Modifier").
        Where("id = ?", id).
        First(&order).Error; err != nil {
        // ...
    }
    return &order, nil
}
```

### 3. GetOrdersByTable

```go
func (r *orderRepository) GetOrdersByTable(tableID uuid.UUID) ([]*models.Order, error) {
    var orders []*models.Order
    if err := r.db.
        Preload("Table").
        Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
            return db.Unscoped() // Include soft-deleted menu items
        }).
        Preload("Items.Modifiers.Modifier").
        Where("table_id = ?", tableID).
        Order("created_at DESC").
        Find(&orders).Error; err != nil {
        // ...
    }
    return orders, nil
}
```

### 4. GetOrdersByStatus

```go
func (r *orderRepository) GetOrdersByStatus(status string) ([]*models.Order, error) {
    var orders []*models.Order
    if err := r.db.
        Preload("Table").
        Preload("Items.MenuItem", func(db *gorm.DB) *gorm.DB {
            return db.Unscoped() // Include soft-deleted menu items
        }).
        Preload("Items.Modifiers.Modifier").
        Where("status = ?", status).
        Order("created_at DESC").
        Find(&orders).Error; err != nil {
        // ...
    }
    return orders, nil
}
```

### 5. GetOrderItemByID

```go
func (r *orderRepository) GetOrderItemByID(id uuid.UUID) (*models.OrderItem, error) {
    var orderItem models.OrderItem
    if err := r.db.
        Preload("MenuItem", func(db *gorm.DB) *gorm.DB {
            return db.Unscoped() // Include soft-deleted menu items
        }).
        Preload("Modifiers.Modifier").
        Where("id = ?", id).
        First(&orderItem).Error; err != nil {
        // ...
    }
    return &orderItem, nil
}
```

## Testing

### Before Fix

```bash
GET /v1/orders
# ❌ Panic: nil pointer dereference
```

### After Fix

```bash
GET /v1/orders
# ✅ Returns all orders with menu item names, even for deleted items
```

## Behavior

### Active Menu Items

```json
{
  "items": [
    {
      "menu_item_id": "abc-123",
      "menu_item_name": "Pad Thai", // Active item - shows normally
      "quantity": 1
    }
  ]
}
```

### Deleted Menu Items

```json
{
  "items": [
    {
      "menu_item_id": "xyz-789",
      "menu_item_name": "Old Special Dish", // Deleted item - still shows in old orders
      "quantity": 1
    }
  ]
}
```

## Important Notes

### ✅ This Is Safe

- Orders only read menu items (no updates)
- Historical data is preserved
- No risk of accidentally "undeleting" items
- Customers see accurate order history

### ⚠️ Don't Use Unscoped() Everywhere

**Use `Unscoped()` only for:**

- Historical/read-only queries (like orders)
- Reports that need complete data
- Audit logs

**Don't use `Unscoped()` for:**

- Menu item list endpoints (customers shouldn't see deleted items)
- Creating new orders (use only active items)
- Menu displays (only show active items)

## Summary

✅ **Fixed:** Orders now load successfully even with deleted menu items  
✅ **Preserved:** Historical order data remains accurate  
✅ **Safe:** Only affects order retrieval, not menu item management  
✅ **Complete:** All order repository methods updated

**Bottom line:** Deleted menu items still appear in historical orders (as they should), but are hidden from active menu lists.
