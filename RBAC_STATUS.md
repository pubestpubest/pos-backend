# RBAC Implementation Status

**Last Updated**: October 16, 2025  
**Status**: 🟡 Foundation Complete - Needs Route Protection

---

## Executive Summary

The POS Backend has **90% of RBAC infrastructure already implemented**. The foundation is solid:

- ✅ Database models with proper relationships
- ✅ Seeded roles and permissions
- ✅ Authentication system with session management
- ✅ Authorization middleware (`RequirePermission`)
- ✅ Permission verification logic

**What's Missing**: Applying the `RequirePermission` middleware to protected routes (2-3 hours of work).

---

## ✅ Completed Components

### 1. Database Schema

- [x] `users` table with UUID primary key
- [x] `roles` table (5 roles: owner, manager, cashier, waiter, kitchen)
- [x] `permissions` table (7 permissions)
- [x] `user_roles` junction table (many-to-many)
- [x] `role_permissions` junction table (many-to-many)
- [x] `sessions` table for token-based authentication

### 2. Domain Layer (`domain/`)

- [x] `AuthUsecase` interface with `VerifyPermission()` and `GetUserPermissions()`
- [x] `AuthRepository` interface with `GetUserPermissions()` and `GetUserWithRolesAndPermissions()`
- [x] `UserUsecase` interface with `AssignRoleToUser()`
- [x] `RoleUsecase` interface
- [x] `PermissionUsecase` interface

### 3. Repository Implementation (`feature/auth/repository/postgres.go`)

- [x] `GetUserPermissions()` - Joins across user_roles, role_permissions, and permissions
- [x] `GetUserWithRolesAndPermissions()` - Eager loads roles and permissions
- [x] `GetSessionByToken()` - Session validation
- [x] Proper error handling with context

### 4. Use Case Implementation (`feature/auth/usecase/usecase.go`)

- [x] `Login()` - Returns permissions in auth response
- [x] `Logout()` - Invalidates session
- [x] `VerifyPermission()` - Checks if user has specific permission
- [x] `GetUserPermissions()` - Returns all user permissions
- [x] `GetUserByToken()` - Validates session and returns user
- [x] Password hashing with bcrypt

### 5. Middleware (`middlewares/auth.go`)

- [x] `AuthMiddleware()` - Validates session token
- [x] `RequirePermission(permissionCode)` - Checks specific permission
- [x] Sets user context for downstream handlers
- ⚠️ **Bug**: Type conversion issue on line 73 (see below)

### 6. Seed Data (`seed/data.go`)

- [x] 5 roles defined
- [x] 7 permissions defined
- [x] Role-permission mappings configured
- [x] Sample users with role assignments
- [x] Seed runner implemented

---

## 🔴 Critical Bug to Fix

**File**: `middlewares/auth.go`  
**Line**: 73  
**Issue**: Type mismatch - userID is stored as string but cast to UUID

### Current Code:

```go
// Line 36 - Sets as string
c.Set("userID", user.ID.String())

// Line 73 - Casts to UUID (will panic!)
hasPermission, err := authUc.VerifyPermission(userID.(uuid.UUID), permissionCode)
```

### Fix Required:

```go
// Option A (Recommended): Store UUID directly
// Change line 36 to:
c.Set("userID", user.ID)

// Keep line 73 as is
hasPermission, err := authUc.VerifyPermission(userID.(uuid.UUID), permissionCode)
```

**Impact**: High - `RequirePermission` middleware will panic until fixed  
**Effort**: 2 minutes  
**Priority**: 🔴 Must fix before applying permissions to routes

---

## 🟡 Pending Implementation

### Route Protection Status

| Route File           | Auth Applied | Permission Applied | Status                     |
| -------------------- | ------------ | ------------------ | -------------------------- |
| `authRoute.go`       | ✅           | N/A                | ✅ Complete                |
| `orderRoute.go`      | ✅ Partial   | ❌                 | 🟡 Needs permission checks |
| `paymentRoute.go`    | ✅ Partial   | ❌                 | 🟡 Needs permission checks |
| `menuItemRoute.go`   | ✅ Partial   | ❌                 | 🟡 Needs permission checks |
| `categoryRoute.go`   | ✅ Partial   | ❌                 | 🟡 Needs permission checks |
| `modifierRoute.go`   | ✅ Partial   | ❌                 | 🟡 Needs permission checks |
| `tableRoute.go`      | ✅           | ❌                 | 🟡 Needs permission checks |
| `areaRoute.go`       | ✅           | ❌                 | 🟡 Needs permission checks |
| `userRoute.go`       | ✅           | ❌                 | 🟡 Needs permission checks |
| `roleRoute.go`       | ✅           | ❌                 | 🟡 Needs permission checks |
| `permissionRoute.go` | ✅           | ❌                 | 🟡 Needs permission checks |

**Note**: ✅ Partial = Some routes protected, some public (by design)

---

## 📋 Implementation Roadmap

### Step 1: Fix Critical Bug (15 minutes)

- [x] Identify the type mismatch issue
- [ ] Change line 36 in `middlewares/auth.go` to store UUID
- [ ] Test `RequirePermission` middleware works
- [ ] Verify no runtime panics

### Step 2: Apply Permission Checks (2-3 hours)

#### High Priority (Core Business Logic)

1. [ ] **orderRoute.go** - Protect close/void operations

   - `PUT /orders/:id/close` → `order.pay`
   - `PUT /orders/:id/void` → `order.pay`

2. [ ] **paymentRoute.go** - Protect financial operations

   - `GET /payments` → `order.pay`
   - `GET /payments/:id` → `order.pay`

3. [ ] **userRoute.go** - Protect user management
   - All endpoints → `user.manage`

#### Medium Priority (Administrative)

4. [ ] **menuItemRoute.go** - Protect menu management

   - POST, PUT, DELETE → `menu.manage`

5. [ ] **categoryRoute.go** - Protect category management

   - POST, PUT, DELETE → `menu.manage`

6. [ ] **modifierRoute.go** - Protect modifier management

   - POST, PUT, DELETE → `menu.manage`

7. [ ] **tableRoute.go** - Protect table operations

   - `PUT /tables/:id/status` → `table.manage`

8. [ ] **areaRoute.go** - Protect area management
   - All endpoints → `table.manage`

#### Low Priority (Read-Only Admin)

9. [ ] **roleRoute.go** - Protect role viewing

   - All endpoints → `user.manage`

10. [ ] **permissionRoute.go** - Protect permission viewing
    - All endpoints → `user.manage`

### Step 3: Testing (2-3 hours)

- [ ] Create test users with different roles
- [ ] Test each protected endpoint with authorized user (should succeed)
- [ ] Test each protected endpoint with unauthorized user (should get 403)
- [ ] Test public endpoints still work without auth
- [ ] Verify error messages are clear

### Step 4: Documentation (1-2 hours)

- [ ] Update API documentation with permission requirements
- [ ] Add error response examples (401, 403)
- [ ] Create Postman collection with examples
- [ ] Update README with RBAC overview

---

## 🎯 Permission Assignment Matrix

### Seeded Roles and Their Permissions

```
owner (All Permissions):
  ✓ order.create
  ✓ order.update
  ✓ order.pay
  ✓ menu.manage
  ✓ table.manage
  ✓ user.manage
  ✓ report.view

manager (All Except User Management):
  ✓ order.create
  ✓ order.update
  ✓ order.pay
  ✓ menu.manage
  ✓ table.manage
  ✗ user.manage
  ✓ report.view

cashier (Payment and Reports):
  ✗ order.create
  ✗ order.update
  ✓ order.pay
  ✗ menu.manage
  ✗ table.manage
  ✗ user.manage
  ✓ report.view

waiter (Order Operations):
  ✓ order.create
  ✓ order.update
  ✗ order.pay
  ✗ menu.manage
  ✗ table.manage
  ✗ user.manage
  ✗ report.view

kitchen (Order Updates Only):
  ✗ order.create
  ✓ order.update
  ✗ order.pay
  ✗ menu.manage
  ✗ table.manage
  ✗ user.manage
  ✗ report.view
```

---

## 📊 Current vs. Target State

### Current State

```
Public Endpoints:
├── /v1/orders/* (customer QR ordering)
├── /v1/menu-items/* (menu viewing)
├── /v1/categories/* (category viewing)
└── /v1/payments/methods (payment methods)

Protected Endpoints (Auth Only):
├── /v1/orders (staff operations) ⚠️ No permission check
├── /v1/payments (payment viewing) ⚠️ No permission check
├── /v1/menu-items (CRUD) ⚠️ No permission check
├── /v1/categories (CRUD) ⚠️ No permission check
├── /v1/modifiers (CRUD) ⚠️ No permission check
├── /v1/tables/* ⚠️ No permission check
├── /v1/areas/* ⚠️ No permission check
├── /v1/users/* ⚠️ No permission check
├── /v1/roles/* ⚠️ No permission check
└── /v1/permissions/* ⚠️ No permission check
```

### Target State

```
Public Endpoints:
├── /v1/orders/* (customer QR ordering)
├── /v1/menu-items/* (menu viewing)
├── /v1/categories/* (category viewing)
└── /v1/payments/methods (payment methods)

Protected Endpoints (With Permission Checks):
├── /v1/orders
│   ├── GET /orders (auth only)
│   ├── GET /orders/open (auth only)
│   ├── PUT /orders/:id/close (order.pay) ✅
│   └── PUT /orders/:id/void (order.pay) ✅
├── /v1/payments
│   ├── GET /payments (order.pay) ✅
│   └── GET /payments/:id (order.pay) ✅
├── /v1/menu-items
│   ├── POST /menu-items (menu.manage) ✅
│   ├── PUT /menu-items/:id (menu.manage) ✅
│   └── DELETE /menu-items/:id (menu.manage) ✅
├── /v1/categories (menu.manage) ✅
├── /v1/modifiers (menu.manage) ✅
├── /v1/tables
│   ├── GET /tables (auth only)
│   └── PUT /tables/:id/status (table.manage) ✅
├── /v1/areas/* (table.manage) ✅
├── /v1/users/* (user.manage) ✅
├── /v1/roles/* (user.manage) ✅
└── /v1/permissions/* (user.manage) ✅
```

---

## 🧪 Testing Scenarios

### Scenario 1: Waiter Login

```bash
# Login as waiter
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "waiter1", "password": "password"}' \
  -c cookies.txt

# Response should include:
{
  "user": {...},
  "token": "...",
  "permissions": ["order.create", "order.update"]
}

# ✅ Should succeed: Create order
curl -X POST http://localhost:8080/v1/orders -b cookies.txt ...

# ❌ Should fail with 403: Close order
curl -X PUT http://localhost:8080/v1/orders/1/close -b cookies.txt
```

### Scenario 2: Cashier Login

```bash
# Login as cashier
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "cashier1", "password": "password"}' \
  -c cookies.txt

# ✅ Should succeed: Close order
curl -X PUT http://localhost:8080/v1/orders/1/close -b cookies.txt

# ❌ Should fail with 403: Modify menu
curl -X POST http://localhost:8080/v1/menu-items -b cookies.txt ...
```

### Scenario 3: Manager Login

```bash
# Login as manager
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "manager1", "password": "password"}' \
  -c cookies.txt

# ✅ Should succeed: Manage menu
curl -X POST http://localhost:8080/v1/menu-items -b cookies.txt ...

# ❌ Should fail with 403: Create user
curl -X POST http://localhost:8080/v1/users -b cookies.txt ...
```

---

## 📁 Related Documentation

1. **RBAC_IMPLEMENTATION_PLAN.md** - Detailed implementation guide with code examples
2. **RBAC_QUICK_REFERENCE.md** - Quick reference for developers
3. **RBAC_STATUS.md** - This file - current status and roadmap

---

## 🚀 Next Actions

### Immediate (Today)

1. Fix the middleware bug (`middlewares/auth.go` line 36)
2. Test the fix with a simple protected endpoint
3. Apply permissions to `userRoute.go` (highest security priority)

### This Week

1. Apply permissions to all route files
2. Test with different user roles
3. Update API documentation

### Future Enhancements

1. Add `RequireAnyPermission()` for OR logic
2. Add `/auth/me` endpoint for current user info
3. Implement permission caching (Redis)
4. Add audit logging for permission denials
5. Create admin UI for role/permission management

---

## ✅ Success Criteria

RBAC implementation is complete when:

- [ ] Middleware bug is fixed and tested
- [ ] All protected routes have appropriate permission checks
- [ ] All 5 roles can only access their permitted endpoints
- [ ] 403 errors are returned for unauthorized access
- [ ] API documentation reflects permission requirements
- [ ] Integration tests cover permission scenarios

---

**Estimated Completion Time**: 1 business day  
**Risk Level**: Low (foundation is solid)  
**Blocker**: None (all dependencies met)

---

## 💡 Quick Start

To implement RBAC right now:

1. **Fix the bug**:

   ```go
   // In middlewares/auth.go, line 36:
   c.Set("userID", user.ID)  // Change from user.ID.String()
   ```

2. **Protect one route** (example):

   ```go
   // In routes/userRoute.go:
   userRoutes.Use(middlewares.AuthMiddleware(), middlewares.RequirePermission("user.manage"))
   ```

3. **Test it**:

   ```bash
   # Login as owner (has user.manage permission)
   # Should succeed

   # Login as waiter (no user.manage permission)
   # Should get 403 Forbidden
   ```

4. **Repeat for all routes** using the permission matrix in RBAC_QUICK_REFERENCE.md

---

**Status**: Ready to implement - all prerequisites met ✅
