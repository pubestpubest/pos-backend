# RBAC Implementation Plan

## Executive Summary

This document outlines the complete implementation plan for Role-Based Access Control (RBAC) in the POS Backend system. The foundation for RBAC already exists - we need to apply permission checks to all protected endpoints.

## Current State Analysis

### ✅ What Already Exists

1. **Database Models** (`models/`)

   - `User` - with many-to-many relationship to Roles
   - `Role` - with many-to-many relationship to Permissions
   - `Permission` - permission codes and descriptions
   - `UserRole` - junction table
   - `RolePermission` - junction table
   - `Session` - token-based session management

2. **Domain Layer** (`domain/`)

   - `AuthUsecase` - includes `VerifyPermission()` and `GetUserPermissions()`
   - `UserRepository` - includes `GetUserWithRoles()`
   - `AuthRepository` - includes `GetUserPermissions()`

3. **Middleware** (`middlewares/auth.go`)

   - `AuthMiddleware()` - validates token and sets user context
   - `RequirePermission(permissionCode)` - checks if user has specific permission

4. **Seed Data** (`seed/data.go`)

   ```
   Roles:
   - owner: Full access (all 7 permissions)
   - manager: All except user.manage
   - cashier: order.pay, report.view
   - waiter: order.create, order.update
   - kitchen: order.update

   Permissions:
   - order.create - Create orders
   - order.update - Update orders
   - order.pay - Take payments
   - menu.manage - CRUD menu & modifiers
   - table.manage - CRUD tables/areas
   - user.manage - Manage users & roles
   - report.view - View reports/dashboard
   ```

### 🔧 What Needs to Be Fixed

1. **Bug in RequirePermission Middleware** (Line 73 in `middlewares/auth.go`)

   - Current: `userID.(uuid.UUID)` - this will panic
   - Issue: userID is stored as string in context (line 36)
   - Fix: Parse string to UUID or store UUID directly

2. **Missing Permission Checks on Routes**
   - Most routes only use `AuthMiddleware()` without permission checks
   - Need to apply `RequirePermission()` to all protected endpoints

## Implementation Plan

### Phase 1: Fix Middleware Bug 🔴 CRITICAL

**File**: `middlewares/auth.go`

**Problem**: Type mismatch in `RequirePermission` middleware

**Current Code** (Lines 63-68):

```go
userID, exists := c.Get("userID")
if !exists {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
    c.Abort()
    return
}
```

**Issue at Line 73**:

```go
hasPermission, err := authUc.VerifyPermission(userID.(uuid.UUID), permissionCode)
```

**Root Cause**: Line 36 sets userID as string: `c.Set("userID", user.ID.String())`

**Solution Options**:

**Option A** (Recommended): Store UUID directly

```go
// Line 36 - Change to store UUID
c.Set("userID", user.ID)

// Line 73 - Keep as is
hasPermission, err := authUc.VerifyPermission(userID.(uuid.UUID), permissionCode)
```

**Option B**: Parse string to UUID

```go
// Line 73 - Parse string
userIDStr := userID.(string)
userUUID, err := uuid.Parse(userIDStr)
if err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
    c.Abort()
    return
}
hasPermission, err := authUc.VerifyPermission(userUUID, permissionCode)
```

**Recommendation**: Use Option A for consistency and performance.

---

### Phase 2: Apply Permission Checks to Routes

Apply `RequirePermission()` middleware to protected endpoints according to the permission matrix below.

#### Permission Matrix

| Endpoint                                  | HTTP Method | Permission Required              | Who Can Access                  |
| ----------------------------------------- | ----------- | -------------------------------- | ------------------------------- |
| **Orders**                                |             |                                  |                                 |
| `POST /orders`                            | POST        | `order.create`                   | Owner, Manager, Waiter          |
| `POST /orders/:id/items`                  | POST        | `order.create`                   | Owner, Manager, Waiter          |
| `PUT /orders/:id/items/:item_id/quantity` | PUT         | `order.update`                   | Owner, Manager, Waiter, Kitchen |
| `DELETE /orders/:id/items/:item_id`       | DELETE      | `order.update`                   | Owner, Manager, Waiter, Kitchen |
| `GET /orders`                             | GET         | `order.pay` or `order.create`    | Owner, Manager, Cashier, Waiter |
| `GET /orders/open`                        | GET         | `order.pay` or `order.create`    | Owner, Manager, Cashier, Waiter |
| `PUT /orders/:id/close`                   | PUT         | `order.pay`                      | Owner, Manager, Cashier         |
| `PUT /orders/:id/void`                    | PUT         | `order.pay`                      | Owner, Manager, Cashier         |
| `GET /tables/:id/orders`                  | GET         | `order.create`                   | Owner, Manager, Waiter          |
| **Payments**                              |             |                                  |                                 |
| `GET /payments`                           | GET         | `order.pay`                      | Owner, Manager, Cashier         |
| `GET /payments/:id`                       | GET         | `order.pay`                      | Owner, Manager, Cashier         |
| **Menu Items**                            |             |                                  |                                 |
| `POST /menu-items`                        | POST        | `menu.manage`                    | Owner, Manager                  |
| `PUT /menu-items/:id`                     | PUT         | `menu.manage`                    | Owner, Manager                  |
| `DELETE /menu-items/:id`                  | DELETE      | `menu.manage`                    | Owner, Manager                  |
| **Categories**                            |             |                                  |                                 |
| `POST /categories`                        | POST        | `menu.manage`                    | Owner, Manager                  |
| `PUT /categories/:id`                     | PUT         | `menu.manage`                    | Owner, Manager                  |
| `DELETE /categories/:id`                  | DELETE      | `menu.manage`                    | Owner, Manager                  |
| **Modifiers**                             |             |                                  |                                 |
| `POST /modifiers`                         | POST        | `menu.manage`                    | Owner, Manager                  |
| `PUT /modifiers/:id`                      | PUT         | `menu.manage`                    | Owner, Manager                  |
| `DELETE /modifiers/:id`                   | DELETE      | `menu.manage`                    | Owner, Manager                  |
| **Tables**                                |             |                                  |                                 |
| `GET /tables`                             | GET         | `table.manage` or `order.create` | Owner, Manager, Waiter          |
| `GET /tables/with-open-orders`            | GET         | `table.manage` or `order.create` | Owner, Manager, Waiter          |
| `GET /tables/:id`                         | GET         | `table.manage` or `order.create` | Owner, Manager, Waiter          |
| `PUT /tables/:id/status`                  | PUT         | `table.manage`                   | Owner, Manager                  |
| **Areas**                                 |             |                                  |                                 |
| `GET /areas`                              | GET         | `table.manage`                   | Owner, Manager                  |
| `GET /areas/with-tables`                  | GET         | `table.manage`                   | Owner, Manager                  |
| `GET /areas/:id`                          | GET         | `table.manage`                   | Owner, Manager                  |
| `POST /areas`                             | POST        | `table.manage`                   | Owner, Manager                  |
| `PUT /areas/:id`                          | PUT         | `table.manage`                   | Owner, Manager                  |
| `DELETE /areas/:id`                       | DELETE      | `table.manage`                   | Owner, Manager                  |
| **Users**                                 |             |                                  |                                 |
| `GET /users`                              | GET         | `user.manage`                    | Owner                           |
| `GET /users/:id`                          | GET         | `user.manage`                    | Owner                           |
| `POST /users`                             | POST        | `user.manage`                    | Owner                           |
| `PUT /users/:id`                          | PUT         | `user.manage`                    | Owner                           |
| `POST /users/:id/roles`                   | POST        | `user.manage`                    | Owner                           |
| **Roles**                                 |             |                                  |                                 |
| `GET /roles`                              | GET         | `user.manage`                    | Owner                           |
| `GET /roles/:id`                          | GET         | `user.manage`                    | Owner                           |
| **Permissions**                           |             |                                  |                                 |
| `GET /permissions`                        | GET         | `user.manage`                    | Owner                           |

---

### Phase 3: Route Implementation Details

#### 3.1 Order Routes (`routes/orderRoute.go`)

**Current Issues**:

- Public routes allow anyone to create/modify orders
- Protected routes only check authentication, not permissions

**Required Changes**:

```go
// Public routes for customers (QR code ordering) - NO CHANGES
orderPublicRoutes := v1.Group("/orders")
{
    orderPublicRoutes.GET("/:id", orderHandler.GetOrderByID)
    orderPublicRoutes.POST("", orderHandler.CreateOrder)
    orderPublicRoutes.POST("/:id/items", orderHandler.AddItemToOrder)
    orderPublicRoutes.DELETE("/:id/items/:item_id", orderHandler.RemoveItemFromOrder)
    orderPublicRoutes.PUT("/:id/items/:item_id/quantity", orderHandler.UpdateOrderItemQuantity)
}

// Protected routes for staff/admin - ADD PERMISSIONS
orderProtectedRoutes := v1.Group("/orders")
orderProtectedRoutes.Use(middlewares.AuthMiddleware())
{
    // View orders - cashier or waiter can view
    orderProtectedRoutes.GET("", orderHandler.GetAllOrders) // No specific permission - any authenticated staff
    orderProtectedRoutes.GET("/open", orderHandler.GetOpenOrders) // No specific permission

    // Close/pay orders - requires order.pay
    orderProtectedRoutes.PUT("/:id/close", middlewares.RequirePermission("order.pay"), orderHandler.CloseOrder)
    orderProtectedRoutes.PUT("/:id/void", middlewares.RequirePermission("order.pay"), orderHandler.VoidOrder)
}

// Table-specific routes
tableOrderRoutes := v1.Group("/tables/:id/orders")
tableOrderRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("order.create"))
{
    tableOrderRoutes.GET("", orderHandler.GetOrdersByTable)
}
```

#### 3.2 Payment Routes (`routes/paymentRoute.go`)

```go
// Public routes - NO CHANGES (for customer self-checkout)
paymentPublicRoutes := v1.Group("/payments")
{
    paymentPublicRoutes.POST("", paymentHandler.ProcessPayment)
    paymentPublicRoutes.GET("/methods", paymentHandler.GetPaymentMethods)
}

// Protected routes - ADD PERMISSIONS
paymentProtectedRoutes := v1.Group("/payments")
paymentProtectedRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("order.pay"))
{
    paymentProtectedRoutes.GET("", paymentHandler.GetAllPayments)
    paymentProtectedRoutes.GET("/:id", paymentHandler.GetPaymentByID)
}

// Order payments - Keep public for customers
orderPaymentPublicRoutes := v1.Group("/orders/:id/payments")
{
    orderPaymentPublicRoutes.GET("", paymentHandler.GetPaymentsByOrder)
}
```

#### 3.3 Menu Item Routes (`routes/menuItemRoute.go`)

```go
// Public routes - NO CHANGES
menuItemPublicRoutes := v1.Group("/menu-items")
{
    menuItemPublicRoutes.GET("", menuItemHandler.GetAllMenuItems)
    menuItemPublicRoutes.GET("/modifiers", menuItemHandler.GetAvailableModifiers)
    menuItemPublicRoutes.GET("/:id", menuItemHandler.GetMenuItemByID)
}

// Protected routes - ADD PERMISSIONS
menuItemProtectedRoutes := v1.Group("/menu-items")
menuItemProtectedRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("menu.manage"))
{
    menuItemProtectedRoutes.POST("", menuItemHandler.CreateMenuItem)
    menuItemProtectedRoutes.PUT("/:id", menuItemHandler.UpdateMenuItem)
    menuItemProtectedRoutes.DELETE("/:id", menuItemHandler.DeleteMenuItem)
}
```

#### 3.4 Category Routes (`routes/categoryRoute.go`)

```go
// Public routes - NO CHANGES
categoryPublicRoutes := v1.Group("/categories")
{
    categoryPublicRoutes.GET("", categoryHandler.GetAllCategories)
    categoryPublicRoutes.GET("/:id", categoryHandler.GetCategoryByID)
}

// Protected routes - ADD PERMISSIONS
categoryProtectedRoutes := v1.Group("/categories")
categoryProtectedRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("menu.manage"))
{
    categoryProtectedRoutes.POST("", categoryHandler.CreateCategory)
    categoryProtectedRoutes.PUT("/:id", categoryHandler.UpdateCategory)
    categoryProtectedRoutes.DELETE("/:id", categoryHandler.DeleteCategory)
}
```

#### 3.5 Modifier Routes (`routes/modifierRoute.go`)

```go
// Public routes - NO CHANGES
modifierPublicRoutes := v1.Group("/modifiers")
{
    modifierPublicRoutes.GET("", modifierHandler.GetAllModifiers)
    modifierPublicRoutes.GET("/:id", modifierHandler.GetModifierByID)
}

categoryModifierRoutes := v1.Group("/categories")
{
    categoryModifierRoutes.GET("/:id/modifiers", modifierHandler.GetModifiersByCategoryID)
}

// Protected routes - ADD PERMISSIONS
modifierProtectedRoutes := v1.Group("/modifiers")
modifierProtectedRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("menu.manage"))
{
    modifierProtectedRoutes.POST("", modifierHandler.CreateModifier)
    modifierProtectedRoutes.PUT("/:id", modifierHandler.UpdateModifier)
    modifierProtectedRoutes.DELETE("/:id", modifierHandler.DeleteModifier)
}
```

#### 3.6 Table Routes (`routes/tableRoute.go`)

```go
tableRoutes := v1.Group("/tables")
tableRoutes.Use(middlewares.AuthMiddleware())
{
    // Read operations - anyone authenticated
    tableRoutes.GET("", tableHandler.GetAllTables)
    tableRoutes.GET("/with-open-orders", tableHandler.GetTablesWithOpenOrders)
    tableRoutes.GET("/:id", tableHandler.GetTableByID)

    // Write operations - requires table.manage
    tableRoutes.PUT("/:id/status", middlewares.RequirePermission("table.manage"), tableHandler.UpdateTableStatus)
}
```

#### 3.7 Area Routes (`routes/areaRoute.go`)

```go
areaRoutes := v1.Group("/areas")
areaRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("table.manage"))
{
    areaRoutes.GET("", areaHandler.GetAllAreas)
    areaRoutes.GET("/with-tables", areaHandler.GetAreasWithTables)
    areaRoutes.GET("/:id", areaHandler.GetAreaByID)
    areaRoutes.POST("", areaHandler.CreateArea)
    areaRoutes.PUT("/:id", areaHandler.UpdateArea)
    areaRoutes.DELETE("/:id", areaHandler.DeleteArea)
}
```

#### 3.8 User Routes (`routes/userRoute.go`)

```go
userRoutes := v1.Group("/users")
userRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("user.manage"))
{
    userRoutes.GET("", userHandler.GetAllUsers)
    userRoutes.GET("/:id", userHandler.GetUserByID)
    userRoutes.POST("", userHandler.CreateUser)
    userRoutes.PUT("/:id", userHandler.UpdateUser)
    userRoutes.POST("/:id/roles", userHandler.AssignRole)
}
```

#### 3.9 Role Routes (`routes/roleRoute.go`)

```go
roleRoutes := v1.Group("/roles")
roleRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("user.manage"))
{
    roleRoutes.GET("", roleHandler.GetAllRoles)
    roleRoutes.GET("/:id", roleHandler.GetRoleWithPermissions)
}
```

#### 3.10 Permission Routes (`routes/permissionRoute.go`)

```go
permissionRoutes := v1.Group("/permissions")
permissionRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("user.manage"))
{
    permissionRoutes.GET("", permissionHandler.GetAllPermissions)
}
```

---

### Phase 4: Enhanced Features (Optional)

#### 4.1 Multiple Permission Support

Create helper middleware for endpoints that accept multiple permissions (OR logic):

```go
// Add to middlewares/auth.go
func RequireAnyPermission(permissionCodes ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID, exists := c.Get("userID")
        if !exists {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }

        authRepo := authRepository.NewAuthRepository(database.DB)
        authUc := authUsecase.NewAuthUsecase(authRepo)

        for _, permissionCode := range permissionCodes {
            hasPermission, err := authUc.VerifyPermission(userID.(uuid.UUID), permissionCode)
            if err == nil && hasPermission {
                c.Next()
                return
            }
        }

        c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        c.Abort()
    }
}
```

**Usage Example**:

```go
// Allow both cashiers and waiters to view orders
tableRoutes.GET("", middlewares.RequireAnyPermission("order.pay", "order.create"), tableHandler.GetAllTables)
```

#### 4.2 Current User Endpoint

Add endpoint for user to get their own info and permissions:

```go
// In routes/authRoute.go
authRoutes.GET("/me", middlewares.AuthMiddleware(), authHandler.GetCurrentUser)
```

#### 4.3 Permission Caching

Implement caching for permission checks to reduce database queries:

```go
// Use Redis or in-memory cache
// Cache key: "user:{userID}:permissions"
// TTL: 5-10 minutes
```

---

## Testing Plan

### 1. Unit Tests

Create tests for:

- `VerifyPermission()` method
- `GetUserPermissions()` method
- `RequirePermission()` middleware
- `RequireAnyPermission()` middleware (if implemented)

### 2. Integration Tests

Test permission enforcement on each protected endpoint:

```go
func TestOrderEndpoints(t *testing.T) {
    tests := []struct {
        name           string
        role           string
        endpoint       string
        method         string
        expectedStatus int
    }{
        {"Owner can close order", "owner", "/v1/orders/1/close", "PUT", 200},
        {"Waiter cannot close order", "waiter", "/v1/orders/1/close", "PUT", 403},
        {"Cashier can close order", "cashier", "/v1/orders/1/close", "PUT", 200},
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### 3. Manual Testing Checklist

Create a test matrix:

| Role    | Action           | Expected Result  |
| ------- | ---------------- | ---------------- |
| Owner   | Create menu item | ✅ Success       |
| Manager | Create menu item | ✅ Success       |
| Cashier | Create menu item | ❌ 403 Forbidden |
| Waiter  | Create order     | ✅ Success       |
| Waiter  | Close order      | ❌ 403 Forbidden |
| Kitchen | Update order     | ✅ Success       |
| Kitchen | View payments    | ❌ 403 Forbidden |

---

## Migration & Deployment

### Pre-Deployment Checklist

- [ ] Database has roles and permissions seeded
- [ ] All users have at least one role assigned
- [ ] Session table exists and is being used
- [ ] All route files updated with permission checks
- [ ] Middleware bug fixed
- [ ] Integration tests passing
- [ ] API documentation updated

### Deployment Steps

1. **Deploy to Staging**

   - Run with `SEED_DB=true` to populate roles/permissions
   - Test all endpoints with different user roles
   - Verify error messages are clear

2. **Monitor**

   - Watch for 403 errors
   - Check logs for permission verification failures
   - Ensure no legitimate users are blocked

3. **Deploy to Production**
   - Seed roles and permissions if fresh install
   - For existing deployments, run migration to add missing data
   - Monitor for 24-48 hours

---

## API Documentation Updates

### Error Responses

Update API documentation to include RBAC errors:

```json
{
  "error": "Insufficient permissions"
}
```

**Status Code**: `403 Forbidden`

### Permission Requirements

Document required permissions for each endpoint:

```markdown
### Close Order

`PUT /v1/orders/:id/close`

**Authentication**: Required  
**Permission**: `order.pay`  
**Roles**: Owner, Manager, Cashier
```

---

## Security Considerations

1. **Session Management**

   - Sessions expire after 24 hours
   - Implement session cleanup cron job
   - Consider refresh token mechanism

2. **Permission Changes**

   - If user's role changes, they need to re-login
   - Consider implementing permission cache invalidation

3. **Audit Logging**

   - Log permission denials
   - Log role assignments/changes
   - Track who modified permissions

4. **Rate Limiting**
   - Implement rate limiting on auth endpoints
   - Prevent brute force attacks

---

## Future Enhancements

1. **Dynamic Permissions**

   - UI for creating custom permissions
   - UI for customizing role permissions

2. **Resource-Level Permissions**

   - Users can only edit their own orders
   - Territory-based permissions (area managers)

3. **Time-Based Permissions**

   - Shift-based access control
   - Scheduled role activation

4. **Permission Inheritance**
   - Role hierarchy (e.g., Manager inherits Waiter permissions)

---

## Summary

### Implementation Checklist

- [ ] **Phase 1**: Fix middleware bug in `middlewares/auth.go`
- [ ] **Phase 2**: Update all route files with permission checks
  - [ ] orderRoute.go
  - [ ] paymentRoute.go
  - [ ] menuItemRoute.go
  - [ ] categoryRoute.go
  - [ ] modifierRoute.go
  - [ ] tableRoute.go
  - [ ] areaRoute.go
  - [ ] userRoute.go
  - [ ] roleRoute.go
  - [ ] permissionRoute.go
- [ ] **Phase 3**: Test all endpoints
- [ ] **Phase 4**: Update API documentation
- [ ] **Phase 5**: Deploy to staging
- [ ] **Phase 6**: Deploy to production

### Estimated Effort

- Middleware fix: **15 minutes**
- Route updates: **2-3 hours**
- Testing: **2-3 hours**
- Documentation: **1-2 hours**
- **Total: 1 day**

---

## Support & Questions

For questions or issues:

1. Review this document
2. Check middleware implementation
3. Verify seed data is loaded
4. Check user role assignments

**End of RBAC Implementation Plan**
