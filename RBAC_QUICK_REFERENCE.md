# RBAC Quick Reference Guide

## Role & Permission Matrix

### Roles

| Role      | Description                           | Use Case                        |
| --------- | ------------------------------------- | ------------------------------- |
| `owner`   | Full system access                    | Restaurant owner/admin          |
| `manager` | All operations except user management | Floor manager, shift supervisor |
| `cashier` | Payment processing and reports        | Cashier, front desk             |
| `waiter`  | Order creation and updates            | Wait staff, servers             |
| `kitchen` | Order updates only                    | Kitchen staff, cooks            |

### Permissions

| Permission Code | Description                            | Roles with Access               |
| --------------- | -------------------------------------- | ------------------------------- |
| `order.create`  | Create new orders                      | owner, manager, waiter          |
| `order.update`  | Modify order items                     | owner, manager, waiter, kitchen |
| `order.pay`     | Process payments and close orders      | owner, manager, cashier         |
| `menu.manage`   | CRUD menu items, categories, modifiers | owner, manager                  |
| `table.manage`  | CRUD tables and areas                  | owner, manager                  |
| `user.manage`   | Manage users and roles                 | owner                           |
| `report.view`   | View reports and dashboard             | owner, manager, cashier         |

---

## Usage in Code

### Protect an Endpoint with Authentication Only

```go
protectedRoutes := v1.Group("/endpoint")
protectedRoutes.Use(middlewares.AuthMiddleware())
{
    protectedRoutes.GET("", handler.GetAll)
}
```

### Protect an Endpoint with Specific Permission

```go
protectedRoutes := v1.Group("/endpoint")
protectedRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("permission.code"))
{
    protectedRoutes.POST("", handler.Create)
}
```

### Apply Permission to Individual Route

```go
protectedRoutes := v1.Group("/endpoint")
protectedRoutes.Use(middlewares.AuthMiddleware())
{
    // Anyone authenticated can read
    protectedRoutes.GET("", handler.GetAll)

    // Only users with permission can write
    protectedRoutes.POST("", middlewares.RequirePermission("permission.code"), handler.Create)
    protectedRoutes.PUT("/:id", middlewares.RequirePermission("permission.code"), handler.Update)
    protectedRoutes.DELETE("/:id", middlewares.RequirePermission("permission.code"), handler.Delete)
}
```

---

## API Response Status Codes

| Status             | Meaning                         | When                                     |
| ------------------ | ------------------------------- | ---------------------------------------- |
| `200 OK`           | Success                         | User has permission                      |
| `401 Unauthorized` | Not authenticated               | Missing or invalid token                 |
| `403 Forbidden`    | Authenticated but no permission | Valid token but insufficient permissions |
| `404 Not Found`    | Resource not found              | -                                        |

---

## Testing with cURL

### 1. Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "owner", "password": "your_password"}' \
  -c cookies.txt
```

### 2. Make Authenticated Request

```bash
curl -X GET http://localhost:8080/v1/users \
  -b cookies.txt
```

### 3. Test Permission (Should succeed for owner)

```bash
curl -X POST http://localhost:8080/v1/users \
  -b cookies.txt \
  -H "Content-Type: application/json" \
  -d '{"username": "newuser", "password": "password123", "full_name": "New User"}'
```

### 4. Test Permission Denial (Login as waiter)

```bash
# Login as waiter
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "waiter1", "password": "your_password"}' \
  -c cookies_waiter.txt

# Try to create user (should get 403)
curl -X POST http://localhost:8080/v1/users \
  -b cookies_waiter.txt \
  -H "Content-Type: application/json" \
  -d '{"username": "newuser", "password": "password123", "full_name": "New User"}'
```

---

## Endpoint Permission Requirements

### Orders

```
GET    /v1/orders                          Auth only
GET    /v1/orders/open                     Auth only
GET    /v1/orders/:id                      Public (customer access)
POST   /v1/orders                          Public (customer access)
POST   /v1/orders/:id/items                Public (customer access)
PUT    /v1/orders/:id/items/:item_id       Public (customer access)
DELETE /v1/orders/:id/items/:item_id       Public (customer access)
PUT    /v1/orders/:id/close                order.pay
PUT    /v1/orders/:id/void                 order.pay
GET    /v1/tables/:id/orders               order.create
```

### Payments

```
GET    /v1/payments                        order.pay
GET    /v1/payments/:id                    order.pay
POST   /v1/payments                        Public (customer access)
GET    /v1/payments/methods                Public
GET    /v1/orders/:id/payments             Public (customer access)
```

### Menu Items

```
GET    /v1/menu-items                      Public
GET    /v1/menu-items/:id                  Public
GET    /v1/menu-items/modifiers            Public
POST   /v1/menu-items                      menu.manage
PUT    /v1/menu-items/:id                  menu.manage
DELETE /v1/menu-items/:id                  menu.manage
```

### Categories

```
GET    /v1/categories                      Public
GET    /v1/categories/:id                  Public
POST   /v1/categories                      menu.manage
PUT    /v1/categories/:id                  menu.manage
DELETE /v1/categories/:id                  menu.manage
```

### Modifiers

```
GET    /v1/modifiers                       Public
GET    /v1/modifiers/:id                   Public
GET    /v1/categories/:id/modifiers        Public
POST   /v1/modifiers                       menu.manage
PUT    /v1/modifiers/:id                   menu.manage
DELETE /v1/modifiers/:id                   menu.manage
```

### Tables

```
GET    /v1/tables                          Auth only
GET    /v1/tables/with-open-orders         Auth only
GET    /v1/tables/:id                      Auth only
PUT    /v1/tables/:id/status               table.manage
```

### Areas

```
GET    /v1/areas                           table.manage
GET    /v1/areas/with-tables               table.manage
GET    /v1/areas/:id                       table.manage
POST   /v1/areas                           table.manage
PUT    /v1/areas/:id                       table.manage
DELETE /v1/areas/:id                       table.manage
```

### Users

```
GET    /v1/users                           user.manage
GET    /v1/users/:id                       user.manage
POST   /v1/users                           user.manage
PUT    /v1/users/:id                       user.manage
POST   /v1/users/:id/roles                 user.manage
```

### Roles & Permissions

```
GET    /v1/roles                           user.manage
GET    /v1/roles/:id                       user.manage
GET    /v1/permissions                     user.manage
```

---

## Common Scenarios

### Scenario 1: Waiter Starting Shift

**Can Do:**

- ✅ View tables and their status
- ✅ Create new orders for tables
- ✅ Add items to orders
- ✅ Update order items
- ✅ View menu items

**Cannot Do:**

- ❌ Close orders / take payment
- ❌ Modify menu items
- ❌ Manage users
- ❌ View financial reports

### Scenario 2: Cashier Working

**Can Do:**

- ✅ Process payments
- ✅ Close orders
- ✅ Void orders
- ✅ View payment history
- ✅ View reports

**Cannot Do:**

- ❌ Create new orders (not their job)
- ❌ Modify menu
- ❌ Manage users

### Scenario 3: Kitchen Staff

**Can Do:**

- ✅ View orders
- ✅ Update order status

**Cannot Do:**

- ❌ Create new orders
- ❌ Take payments
- ❌ Modify menu
- ❌ Manage anything else

### Scenario 4: Manager Duties

**Can Do:**

- ✅ Everything except user management
- ✅ Create/update orders
- ✅ Process payments
- ✅ Manage menu items
- ✅ Manage tables and areas
- ✅ View reports

**Cannot Do:**

- ❌ Create/delete users
- ❌ Assign roles to users

### Scenario 5: Owner Access

**Can Do:**

- ✅ **Everything** - full system access

---

## Troubleshooting

### Getting 401 Unauthorized

**Problem**: Not logged in or session expired  
**Solution**: Login again or check if token cookie is being sent

### Getting 403 Forbidden

**Problem**: Logged in but insufficient permissions  
**Solution**:

1. Check which permission is required for the endpoint
2. Verify user has the correct role assigned
3. Verify role has the required permission

### Permission Check Not Working

**Checklist**:

1. ✅ Is `AuthMiddleware()` applied before `RequirePermission()`?
2. ✅ Is the permission code spelled correctly?
3. ✅ Does the role have this permission in seed data?
4. ✅ Has the user been assigned the role?
5. ✅ Is the middleware bug fixed? (userID type issue)

---

## Database Queries for Debugging

### Check User's Roles

```sql
SELECT u.username, r.name as role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE u.username = 'waiter1';
```

### Check Role's Permissions

```sql
SELECT r.name as role_name, p.code as permission_code, p.description
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE r.name = 'manager';
```

### Check User's Effective Permissions

```sql
SELECT DISTINCT p.code, p.description
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
WHERE u.username = 'manager1';
```

### Check All Seeded Data

```sql
-- Roles
SELECT * FROM roles;

-- Permissions
SELECT * FROM permissions;

-- Role-Permission mapping
SELECT r.name as role, p.code as permission
FROM roles r
JOIN role_permissions rp ON r.id = rp.role_id
JOIN permissions p ON rp.permission_id = p.id
ORDER BY r.name, p.code;
```

---

## Quick Implementation Checklist

When adding RBAC to a new endpoint:

- [ ] Determine if endpoint should be public or protected
- [ ] If protected, apply `middlewares.AuthMiddleware()`
- [ ] Determine required permission from the matrix above
- [ ] Apply `middlewares.RequirePermission("permission.code")`
- [ ] Test with user that has permission (should succeed)
- [ ] Test with user without permission (should get 403)
- [ ] Update API documentation

---

## Example: Adding New Protected Endpoint

```go
// Step 1: Define the route group
productRoutes := v1.Group("/products")

// Step 2: Apply authentication middleware
productRoutes.Use(middlewares.AuthMiddleware())

// Step 3: Define routes with appropriate permissions
{
    // Public read for authenticated users
    productRoutes.GET("", productHandler.GetAll)
    productRoutes.GET("/:id", productHandler.GetByID)

    // Write operations require menu.manage permission
    productRoutes.POST("",
        middlewares.RequirePermission("menu.manage"),
        productHandler.Create)

    productRoutes.PUT("/:id",
        middlewares.RequirePermission("menu.manage"),
        productHandler.Update)

    productRoutes.DELETE("/:id",
        middlewares.RequirePermission("menu.manage"),
        productHandler.Delete)
}
```

---

**Last Updated**: October 2025  
**Version**: 1.0  
**Related**: See `RBAC_IMPLEMENTATION_PLAN.md` for detailed implementation guide
