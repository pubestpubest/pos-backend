# Postman Testing Guide - Order API Endpoints

## 📋 Table of Contents

1. [Prerequisites](#prerequisites)
2. [Environment Setup](#environment-setup)
3. [Authentication](#authentication)
4. [Order Endpoints](#order-endpoints)
5. [Complete Testing Workflow](#complete-testing-workflow)
6. [Common Errors](#common-errors)

---

## 🚀 Prerequisites

### Required Setup

- Server running on `http://localhost:8080`
- Database seeded with test data (set `SEED_DB=true` in `.env`)
- Postman installed
- Valid user account for authentication

### Test Data Requirements

Before testing orders, ensure you have:

- ✅ At least one user account
- ✅ At least one area
- ✅ At least one table (linked to an area)
- ✅ At least one category
- ✅ At least one menu item (linked to a category)
- ✅ (Optional) Some modifiers for menu items

---

## ⚙️ Environment Setup

### Create Postman Environment

1. Click on "Environments" in Postman
2. Create new environment: `POS Backend - Local`
3. Add variables:

| Variable     | Initial Value              | Current Value                      |
| ------------ | -------------------------- | ---------------------------------- |
| base_url     | http://localhost:8080      | http://localhost:8080              |
| api_version  | v1                         | v1                                 |
| token        | (empty)                    | (will be set after login)          |
| order_id     | (empty)                    | (will be set after creating order) |
| item_id      | (empty)                    | (will be set after adding item)    |
| table_id     | (get from GET /tables)     | (UUID)                             |
| menu_item_id | (get from GET /menu-items) | (UUID)                             |

### Postman Collection Structure

```
POS Backend
├── Auth
│   └── Login
├── Setup (Get Test Data IDs)
│   ├── Get All Tables
│   └── Get All Menu Items
└── Orders
    ├── Create Order
    ├── Get All Orders
    ├── Get Open Orders
    ├── Get Order By ID
    ├── Add Item to Order
    ├── Update Item Quantity
    ├── Remove Item from Order
    ├── Close Order
    └── Void Order
```

---

## 🔐 Authentication

### 1. Login to Get Token

**Endpoint:** `POST {{base_url}}/{{api_version}}/auth/login`

**Headers:**

```
Content-Type: application/json
```

**Request Body:**

```json
{
  "username": "admin",
  "password": "admin123"
}
```

**Expected Response (200 OK):**

```json
{
  "token": "550e8400-e29b-41d4-a716-446655440000",
  "user": {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "username": "admin",
    "full_name": "Administrator",
    "email": "admin@pos.com",
    "phone": "1234567890",
    "status": "active"
  }
}
```

**Postman Test Script (Tests tab):**

```javascript
if (pm.response.code === 200) {
  const response = pm.response.json();
  pm.environment.set("token", response.token);
  console.log("Token saved:", response.token);
}
```

### 2. Set Authorization for All Order Requests

For ALL subsequent requests, add to **Headers:**

```
Authorization: Bearer {{token}}
```

Or use Postman's Authorization tab:

- Type: `Bearer Token`
- Token: `{{token}}`

---

## 📦 Order Endpoints

### 1. Create Order

**Endpoint:** `POST {{base_url}}/{{api_version}}/orders`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {{token}}
```

**Request Body:**

```json
{
  "table_id": "{{table_id}}",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "note": "Customer prefers window seat"
}
```

**Field Descriptions:**

- `table_id` (required): UUID of the table
- `opened_by` (required): UUID of the user creating the order
- `source` (required): Either `"staff"` or `"customer"`
- `note` (optional): Additional notes about the order

**Expected Response (201 Created):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "table_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "Table 1",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "status": "open",
  "subtotal_baht": 0,
  "discount_baht": 0,
  "total_baht": 0,
  "note": "Customer prefers window seat",
  "created_at": "2025-10-12T10:30:00Z",
  "closed_at": null,
  "items": []
}
```

**Postman Test Script:**

```javascript
if (pm.response.code === 201) {
  const response = pm.response.json();
  pm.environment.set("order_id", response.id);
  console.log("Order created:", response.id);
}
```

**What to Expect:**

- ✅ Order is created with status "open"
- ✅ Total amounts are 0 (no items yet)
- ✅ `created_at` is set to current time
- ✅ `closed_at` is null (order is open)
- ✅ Items array is empty

---

### 2. Get All Orders

**Endpoint:** `GET {{base_url}}/{{api_version}}/orders`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
[
  {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "table_id": "550e8400-e29b-41d4-a716-446655440000",
    "table_name": "Table 1",
    "opened_by": "123e4567-e89b-12d3-a456-426614174000",
    "source": "staff",
    "status": "open",
    "subtotal_baht": 15000,
    "discount_baht": 0,
    "total_baht": 15000,
    "note": "Customer prefers window seat",
    "created_at": "2025-10-12T10:30:00Z",
    "closed_at": null,
    "items": [...]
  },
  {
    "id": "8d9e7680-8425-50de-955b-f08fd2f01bf8",
    "table_id": "660f9500-f30c-52e5-b827-557766551111",
    "table_name": "Table 2",
    "opened_by": "123e4567-e89b-12d3-a456-426614174000",
    "source": "customer",
    "status": "paid",
    "subtotal_baht": 25000,
    "discount_baht": 0,
    "total_baht": 25000,
    "note": "",
    "created_at": "2025-10-12T09:15:00Z",
    "closed_at": "2025-10-12T10:00:00Z",
    "items": [...]
  }
]
```

**What to Expect:**

- ✅ Returns array of all orders
- ✅ Includes orders with all statuses (open, paid, void)
- ✅ Each order includes basic info + items array
- ✅ Amounts are in Satang (1 Baht = 100 Satang, so 15000 = 150 Baht)

---

### 3. Get Open Orders Only

**Endpoint:** `GET {{base_url}}/{{api_version}}/orders/open`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
[
  {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "table_id": "550e8400-e29b-41d4-a716-446655440000",
    "table_name": "Table 1",
    "opened_by": "123e4567-e89b-12d3-a456-426614174000",
    "source": "staff",
    "status": "open",
    "subtotal_baht": 15000,
    "discount_baht": 0,
    "total_baht": 15000,
    "note": "Customer prefers window seat",
    "created_at": "2025-10-12T10:30:00Z",
    "closed_at": null,
    "items": [...]
  }
]
```

**What to Expect:**

- ✅ Returns only orders with status = "open"
- ✅ Useful for showing active orders in the system
- ✅ `closed_at` will always be null for these orders

---

### 4. Get Order By ID

**Endpoint:** `GET {{base_url}}/{{api_version}}/orders/{{order_id}}`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "table_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "Table 1",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "status": "open",
  "subtotal_baht": 15000,
  "discount_baht": 0,
  "total_baht": 15000,
  "note": "Customer prefers window seat",
  "created_at": "2025-10-12T10:30:00Z",
  "closed_at": null,
  "items": [
    {
      "id": "9e0f7781-9536-61ef-a66c-g19ge3g12cg9",
      "menu_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "menu_item_name": "Pad Thai",
      "quantity": 2,
      "unit_price_baht": 5000,
      "line_total_baht": 10000,
      "note": "Extra spicy",
      "modifiers": [
        {
          "modifier_id": "b2c3d4e5-f6g7-8901-bcde-fg2345678901",
          "modifier_name": "Extra Shrimp",
          "price_delta_baht": 2000
        }
      ]
    },
    {
      "id": "0f1g8892-a647-72fg-b77d-h20hf4h23dh0",
      "menu_item_id": "c3d4e5f6-g7h8-9012-cdef-gh3456789012",
      "menu_item_name": "Thai Iced Tea",
      "quantity": 1,
      "unit_price_baht": 5000,
      "line_total_baht": 5000,
      "note": "",
      "modifiers": []
    }
  ]
}
```

**What to Expect:**

- ✅ Returns single order with full details
- ✅ Includes all items with modifiers
- ✅ Prices calculated correctly:
  - Unit price: base menu item price
  - Line total: (unit_price + modifier prices) × quantity
  - Order total: sum of all line totals

**Error Response (404 Not Found):**

```json
{
  "error": "Order not found"
}
```

---

### 5. Add Item to Order

**Endpoint:** `POST {{base_url}}/{{api_version}}/orders/{{order_id}}/items`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {{token}}
```

**Request Body (Without Modifiers):**

```json
{
  "menu_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "quantity": 2,
  "note": "Extra spicy"
}
```

**Request Body (With Modifiers):**

```json
{
  "menu_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "quantity": 2,
  "note": "Extra spicy",
  "modifier_ids": [
    "b2c3d4e5-f6g7-8901-bcde-fg2345678901",
    "c3d4e5f6-g7h8-9012-cdef-gh3456789012"
  ]
}
```

**Field Descriptions:**

- `menu_item_id` (required): UUID of the menu item to add
- `quantity` (required): Number of items (minimum: 1)
- `note` (optional): Special instructions for this item
- `modifier_ids` (optional): Array of modifier UUIDs to apply

**Expected Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "table_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "Table 1",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "status": "open",
  "subtotal_baht": 10000,
  "discount_baht": 0,
  "total_baht": 10000,
  "note": "Customer prefers window seat",
  "created_at": "2025-10-12T10:30:00Z",
  "closed_at": null,
  "items": [
    {
      "id": "9e0f7781-9536-61ef-a66c-g19ge3g12cg9",
      "menu_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "menu_item_name": "Pad Thai",
      "quantity": 2,
      "unit_price_baht": 5000,
      "line_total_baht": 10000,
      "note": "Extra spicy",
      "modifiers": []
    }
  ]
}
```

**Postman Test Script:**

```javascript
if (pm.response.code === 200) {
  const response = pm.response.json();
  if (response.items && response.items.length > 0) {
    // Save the first item ID for later use
    pm.environment.set("item_id", response.items[0].id);
    console.log("Item added:", response.items[0].id);
  }
}
```

**What to Expect:**

- ✅ Item is added to the order
- ✅ Order total is updated automatically
- ✅ Returns complete updated order with all items
- ✅ If modifiers specified, they appear in the response
- ✅ Line total = (menu_item_price + sum of modifier prices) × quantity

---

### 6. Update Order Item Quantity

**Endpoint:** `PUT {{base_url}}/{{api_version}}/orders/{{order_id}}/items/{{item_id}}/quantity`

**Headers:**

```
Content-Type: application/json
Authorization: Bearer {{token}}
```

**Request Body:**

```json
{
  "quantity": 3
}
```

**Field Descriptions:**

- `quantity` (required): New quantity (minimum: 1)

**Expected Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "table_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "Table 1",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "status": "open",
  "subtotal_baht": 15000,
  "discount_baht": 0,
  "total_baht": 15000,
  "note": "Customer prefers window seat",
  "created_at": "2025-10-12T10:30:00Z",
  "closed_at": null,
  "items": [
    {
      "id": "9e0f7781-9536-61ef-a66c-g19ge3g12cg9",
      "menu_item_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
      "menu_item_name": "Pad Thai",
      "quantity": 3,
      "unit_price_baht": 5000,
      "line_total_baht": 15000,
      "note": "Extra spicy",
      "modifiers": []
    }
  ]
}
```

**What to Expect:**

- ✅ Item quantity is updated
- ✅ Line total recalculated: unit_price × new_quantity
- ✅ Order total updated automatically
- ✅ Returns complete updated order

**Validation:**

- ❌ Quantity must be at least 1
- ❌ If you want to remove, use Delete endpoint instead

---

### 7. Remove Item from Order

**Endpoint:** `DELETE {{base_url}}/{{api_version}}/orders/{{order_id}}/items/{{item_id}}`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "table_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "Table 1",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "status": "open",
  "subtotal_baht": 0,
  "discount_baht": 0,
  "total_baht": 0,
  "note": "Customer prefers window seat",
  "created_at": "2025-10-12T10:30:00Z",
  "closed_at": null,
  "items": []
}
```

**What to Expect:**

- ✅ Item is removed from order
- ✅ Order total updated automatically
- ✅ Returns complete updated order
- ✅ If it was the last item, items array will be empty

---

### 8. Get Orders by Table

**Endpoint:** `GET {{base_url}}/{{api_version}}/tables/{{table_id}}/orders`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
[
  {
    "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
    "table_id": "550e8400-e29b-41d4-a716-446655440000",
    "table_name": "Table 1",
    "opened_by": "123e4567-e89b-12d3-a456-426614174000",
    "source": "staff",
    "status": "open",
    "subtotal_baht": 15000,
    "discount_baht": 0,
    "total_baht": 15000,
    "note": "Customer prefers window seat",
    "created_at": "2025-10-12T10:30:00Z",
    "closed_at": null,
    "items": [...]
  }
]
```

**What to Expect:**

- ✅ Returns all orders for specified table
- ✅ Includes orders with all statuses
- ✅ Useful for viewing table history
- ✅ Empty array if table has no orders

---

### 9. Close Order

**Endpoint:** `PUT {{base_url}}/{{api_version}}/orders/{{order_id}}/close`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
{
  "id": "7c9e6679-7425-40de-944b-e07fc1f90ae7",
  "table_id": "550e8400-e29b-41d4-a716-446655440000",
  "table_name": "Table 1",
  "opened_by": "123e4567-e89b-12d3-a456-426614174000",
  "source": "staff",
  "status": "paid",
  "subtotal_baht": 15000,
  "discount_baht": 0,
  "total_baht": 15000,
  "note": "Customer prefers window seat",
  "created_at": "2025-10-12T10:30:00Z",
  "closed_at": "2025-10-12T11:00:00Z",
  "items": [...]
}
```

**What to Expect:**

- ✅ Order status changes from "open" to "paid"
- ✅ `closed_at` timestamp is set to current time
- ✅ Order can no longer be modified after closing
- ❌ Cannot close an already closed order
- ❌ Cannot close a voided order

---

### 10. Void Order

**Endpoint:** `PUT {{base_url}}/{{api_version}}/orders/{{order_id}}/void`

**Headers:**

```
Authorization: Bearer {{token}}
```

**Request Body:** None

**Expected Response (200 OK):**

```json
{
  "message": "Order voided successfully"
}
```

**What to Expect:**

- ✅ Order status changes to "void"
- ✅ Order is cancelled and cannot be modified
- ✅ Useful for cancelled orders or mistakes
- ❌ Cannot void an already paid order
- ❌ Cannot void an already voided order

---

## 🔄 Complete Testing Workflow

### Scenario: Create and Process a Complete Order

Follow these steps in order:

#### Step 1: Login

```
POST /v1/auth/login
```

→ Save the `token`

#### Step 2: Get Table ID

```
GET /v1/tables
```

→ Choose a table, save the `table_id`

#### Step 3: Get Menu Item IDs

```
GET /v1/menu-items
```

→ Choose menu items, save their IDs

#### Step 4: Create Order

```
POST /v1/orders
Body: {
  "table_id": "{{table_id}}",
  "opened_by": "{{user_id}}",
  "source": "staff"
}
```

→ Save the `order_id`
→ **Expected**: Order created with status "open", total = 0

#### Step 5: Add First Item

```
POST /v1/orders/{{order_id}}/items
Body: {
  "menu_item_id": "{{menu_item_id_1}}",
  "quantity": 2,
  "note": "Extra spicy"
}
```

→ Save the first `item_id`
→ **Expected**: Order total updated, item appears in order

#### Step 6: Add Second Item with Modifiers

```
POST /v1/orders/{{order_id}}/items
Body: {
  "menu_item_id": "{{menu_item_id_2}}",
  "quantity": 1,
  "modifier_ids": ["{{modifier_id}}"]
}
```

→ **Expected**: Order total increases, item with modifiers added

#### Step 7: Update Item Quantity

```
PUT /v1/orders/{{order_id}}/items/{{item_id}}/quantity
Body: {
  "quantity": 3
}
```

→ **Expected**: Item quantity updated, totals recalculated

#### Step 8: View Current Order

```
GET /v1/orders/{{order_id}}
```

→ **Expected**: See complete order with all items and correct totals

#### Step 9: Remove an Item

```
DELETE /v1/orders/{{order_id}}/items/{{item_id}}
```

→ **Expected**: Item removed, totals updated

#### Step 10: View Open Orders

```
GET /v1/orders/open
```

→ **Expected**: Your order appears in the list

#### Step 11: Close Order

```
PUT /v1/orders/{{order_id}}/close
```

→ **Expected**: Order status = "paid", closed_at timestamp set

#### Step 12: Verify Closed Order

```
GET /v1/orders/{{order_id}}
```

→ **Expected**: Order shows as paid with closed_at timestamp

---

## ⚠️ Common Errors

### Authentication Errors

**Error:** `Authorization header required`

```json
{
  "error": "Authorization header required"
}
```

**Solution:** Add `Authorization: Bearer {{token}}` header

**Error:** `Invalid or expired token`

```json
{
  "error": "Invalid or expired token"
}
```

**Solution:** Login again to get a new token

---

### Validation Errors

**Error:** `Invalid order ID`

```json
{
  "error": "Invalid order ID"
}
```

**Cause:** Order ID is not a valid UUID
**Solution:** Check the UUID format, ensure it's correct

**Error:** `Invalid request body`

```json
{
  "error": "Invalid request body"
}
```

**Cause:** Missing required fields or invalid JSON
**Solution:** Check request body matches the required format

**Error:** Quantity validation

```json
{
  "error": "Invalid request body"
}
```

**Cause:** Quantity is less than 1 or not provided
**Solution:** Set quantity to at least 1

**Error:** Source validation

```json
{
  "error": "Invalid request body"
}
```

**Cause:** Source field is not "staff" or "customer"
**Solution:** Use only "staff" or "customer" for source field

---

### Business Logic Errors

**Error:** `Error getting order`

```json
{
  "error": "Error getting order"
}
```

**Cause:** Order not found with given ID
**Solution:** Verify order ID exists in database

**Error:** `Error creating order`

```json
{
  "error": "Error creating order"
}
```

**Possible Causes:**

- Table ID doesn't exist
- Table is already occupied (depends on business logic)
- User ID doesn't exist

**Error:** `Error adding item to order`

```json
{
  "error": "Error adding item to order"
}
```

**Possible Causes:**

- Menu item ID doesn't exist
- Modifier ID doesn't exist
- Order is already closed/voided

**Error:** `Error closing order`

```json
{
  "error": "Error closing order"
}
```

**Possible Causes:**

- Order is already closed
- Order is already voided
- Order doesn't exist

---

## 📊 Price Calculations Explained

### Understanding Baht/Satang

- Prices are stored in **Satang** (smallest Thai currency unit)
- 1 Baht = 100 Satang
- Example: 15000 Satang = 150 Baht

### Order Total Calculation

```
For each item:
  Base Price = menu_item.price
  Modifiers Total = sum(modifier.price for each modifier)
  Unit Price = Base Price + Modifiers Total
  Line Total = Unit Price × Quantity

Order Subtotal = sum(Line Total for all items)
Order Discount = 0 (if no discount applied)
Order Total = Subtotal - Discount
```

### Example Calculation

**Menu Item:** Pad Thai = 50 Baht (5000 Satang)
**Modifier:** Extra Shrimp = 20 Baht (2000 Satang)
**Quantity:** 2

```
Unit Price = 5000 + 2000 = 7000 Satang (70 Baht)
Line Total = 7000 × 2 = 14000 Satang (140 Baht)
```

---

## 📌 Testing Tips

### 1. Save Response IDs

Use Postman's **Test Scripts** to automatically save IDs:

```javascript
// After creating order
if (pm.response.code === 201) {
  const response = pm.response.json();
  pm.environment.set("order_id", response.id);
}
```

### 2. Create a Test Data Collection

Before testing orders, create:

1. Test user account
2. Test area
3. Test table
4. Test category
5. Test menu items
6. Test modifiers

### 3. Use Postman Collection Runner

- Organize requests in logical order
- Use Collection Runner to execute complete workflows
- Add delays between requests if needed

### 4. Check Database State

After operations, verify in database:

```sql
-- Check order status
SELECT * FROM orders WHERE id = 'order-uuid';

-- Check order items
SELECT * FROM order_items WHERE order_id = 'order-uuid';

-- Check order totals
SELECT
  o.id,
  o.status,
  SUM(oi.quantity * (mi.price)) as calculated_total,
  o.total_amount as stored_total
FROM orders o
LEFT JOIN order_items oi ON oi.order_id = o.id
LEFT JOIN menu_items mi ON mi.id = oi.menu_item_id
WHERE o.id = 'order-uuid'
GROUP BY o.id;
```

### 5. Test Edge Cases

- Create order with no items
- Add same item multiple times
- Update quantity to very large numbers
- Try to modify closed orders
- Try to close orders with no items

---

## 🎯 Quick Reference

### Order Status Flow

```
[CREATE] → open → [CLOSE] → paid
               ↓
              [VOID] → void
```

### HTTP Status Codes

- `200 OK` - Successful GET, PUT, DELETE
- `201 Created` - Successful POST (create)
- `400 Bad Request` - Invalid input/validation error
- `401 Unauthorized` - Missing or invalid token
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server/database error

### Required Headers

All requests (except login):

```
Authorization: Bearer {{token}}
```

All POST/PUT requests:

```
Content-Type: application/json
```

---

**Document Version:** 1.0  
**Last Updated:** October 12, 2025  
**Base URL:** `http://localhost:8080/v1`
