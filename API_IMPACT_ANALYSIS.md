# API Impact Analysis: MenuItem Image Upload Changes

## Overview

The MenuItem endpoints have been changed from **JSON** to **multipart/form-data** to support image uploads. This is a **BREAKING CHANGE** for MenuItem Create and Update operations.

---

## 🔴 Affected Endpoints (BREAKING CHANGES)

### 1. Create Menu Item

**Endpoint:** `POST /v1/menu-items`

**Status:** ⚠️ **BREAKING CHANGE**

**Before:**

- Content-Type: `application/json`
- Body format: JSON

**After:**

- Content-Type: `multipart/form-data`
- Body format: Form data with file upload

**Impact:**

- ❌ Old JSON requests will FAIL
- ✅ Must update to use FormData
- ✅ Can now upload images
- ✅ Image upload is optional

**Migration Required:** YES

---

### 2. Update Menu Item

**Endpoint:** `PUT /v1/menu-items/:id`

**Status:** ⚠️ **BREAKING CHANGE**

**Before:**

- Content-Type: `application/json`
- Body format: JSON

**After:**

- Content-Type: `multipart/form-data`
- Body format: Form data with file upload

**Impact:**

- ❌ Old JSON requests will FAIL
- ✅ Must update to use FormData
- ✅ Can now upload/replace images
- ✅ Omitting image field preserves existing image
- ✅ Automatically deletes old image when uploading new one

**Migration Required:** YES

---

### 3. Delete Menu Item

**Endpoint:** `DELETE /v1/menu-items/:id`

**Status:** ✅ **NO BREAKING CHANGE** (Enhanced - Soft Delete)

**Changes:**

- Same API signature
- Now uses **soft delete** (sets `deleted_at` timestamp instead of physical deletion)
- No longer throws foreign key errors when item is in orders
- Image is kept in MinIO (can be restored if needed)

**Impact:**

- ✅ No client changes required
- ✅ Works even if item is referenced in orders
- ✅ Better data preservation and audit trail
- ⚠️ Database migration required (add `deleted_at` column)

**Migration Required:** YES (Database only - see migrations/add_soft_delete_to_menu_items.sql)

---

## ✅ Unaffected Endpoints (NO CHANGES)

The following endpoints continue to work exactly as before:

### Menu Items (Read Operations)

#### Get All Menu Items

**Endpoint:** `GET /v1/menu-items`

- Status: ✅ No changes
- Still returns JSON with `image_url` field

#### Get Menu Item by ID

**Endpoint:** `GET /v1/menu-items/:id`

- Status: ✅ No changes
- Still returns JSON with `image_url` field

#### Get Available Modifiers

**Endpoint:** `GET /v1/menu-items/modifiers`

- Status: ✅ No changes

---

### All Other Features (Unchanged)

These features continue to use **JSON** format:

#### Categories

- `GET /v1/categories` ✅
- `GET /v1/categories/:id` ✅
- `POST /v1/categories` ✅
- `PUT /v1/categories/:id` ✅
- `DELETE /v1/categories/:id` ✅

#### Modifiers

- `GET /v1/modifiers` ✅
- `GET /v1/modifiers/:id` ✅
- `POST /v1/modifiers` ✅
- `PUT /v1/modifiers/:id` ✅
- `DELETE /v1/modifiers/:id` ✅

#### Areas

- `GET /v1/areas` ✅
- `GET /v1/areas/:id` ✅
- `POST /v1/areas` ✅
- `PUT /v1/areas/:id` ✅
- `DELETE /v1/areas/:id` ✅

#### Tables

- `GET /v1/tables` ✅
- `GET /v1/tables/:id` ✅
- `POST /v1/tables` ✅
- `PUT /v1/tables/:id` ✅
- `DELETE /v1/tables/:id` ✅

#### Orders

- `GET /v1/orders` ✅
- `GET /v1/orders/:id` ✅
- `POST /v1/orders` ✅
- `PUT /v1/orders/:id` ✅
- `DELETE /v1/orders/:id` ✅

#### Payments

- `GET /v1/payments` ✅
- `POST /v1/payments` ✅

#### Auth

- `POST /v1/auth/login` ✅
- `POST /v1/auth/logout` ✅
- `POST /v1/auth/refresh` ✅

#### Users

- `GET /v1/users` ✅
- `GET /v1/users/:id` ✅
- `POST /v1/users` ✅
- `PUT /v1/users/:id` ✅
- `DELETE /v1/users/:id` ✅

#### Roles

- `GET /v1/roles` ✅
- `GET /v1/roles/:id` ✅
- `POST /v1/roles` ✅
- `PUT /v1/roles/:id` ✅
- `DELETE /v1/roles/:id` ✅

#### Permissions

- `GET /v1/permissions` ✅
- `GET /v1/permissions/:id` ✅

---

## Summary Table

| Endpoint                   | Method | Status    | Breaking? | Migration Required? |
| -------------------------- | ------ | --------- | --------- | ------------------- |
| `/v1/menu-items`           | POST   | Changed   | ✅ YES    | ✅ YES              |
| `/v1/menu-items/:id`       | PUT    | Changed   | ✅ YES    | ✅ YES              |
| `/v1/menu-items/:id`       | DELETE | Enhanced  | ❌ NO     | ❌ NO               |
| `/v1/menu-items`           | GET    | Unchanged | ❌ NO     | ❌ NO               |
| `/v1/menu-items/:id`       | GET    | Unchanged | ❌ NO     | ❌ NO               |
| `/v1/menu-items/modifiers` | GET    | Unchanged | ❌ NO     | ❌ NO               |
| All other endpoints        | ALL    | Unchanged | ❌ NO     | ❌ NO               |

---

## Breaking Change Details

### What Will Break?

#### 1. Existing Frontend Code

```typescript
// ❌ THIS WILL FAIL NOW
const response = await fetch("/v1/menu-items", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
  },
  body: JSON.stringify({
    name: "Item",
    sku: "SKU123",
    price_baht: 100,
  }),
});
```

**Error:** `400 Bad Request - Failed to parse form data`

#### 2. Postman Collections

- All saved POST/PUT requests for menu items will fail
- Need to change from JSON to form-data

#### 3. Third-party Integrations

- Any external services creating/updating menu items will break
- Need to update their integration code

---

## Response Format Changes

### Menu Item Responses (Enhanced)

**Before:**

```json
{
  "id": "uuid",
  "name": "ข้าวหมูกรอบพริกเกลือ",
  "sku": "MOU-ASD-PRI-3",
  "price_baht": 70,
  "active": true,
  "image_url": "",  // Usually empty or external URL
  "category": { ... }
}
```

**After:**

```json
{
  "id": "uuid",
  "name": "ข้าวหมูกรอบพริกเกลือ",
  "sku": "MOU-ASD-PRI-3",
  "price_baht": 70,
  "active": true,
  "image_url": "http://localhost:9000/menu-images/uuid-timestamp.jpg",  // MinIO URL
  "category": { ... }
}
```

**Changes:**

- ✅ Response structure is the same
- ✅ `image_url` now contains MinIO URLs for uploaded images
- ✅ Backward compatible for read operations

---

## Migration Checklist

### High Priority (Breaking Changes)

- [ ] **Update Frontend - Menu Item Create Form**

  - Change from JSON to FormData
  - Add file input for image upload
  - Update fetch/axios calls

- [ ] **Update Frontend - Menu Item Update Form**

  - Change from JSON to FormData
  - Add file input for image upload
  - Update fetch/axios calls

- [ ] **Update Postman Collection**

  - Change POST /v1/menu-items to form-data
  - Change PUT /v1/menu-items/:id to form-data
  - Update field types

- [ ] **Notify Third-party Integrations**
  - Email partners about API changes
  - Provide migration guide
  - Set migration deadline

### Medium Priority (Environment Setup)

- [ ] **Add MinIO Environment Variables**

  - Add to `configs/.env`
  - Update deployment configs

- [ ] **Start MinIO Service**

  - Run `docker-compose up -d`
  - Verify MinIO is accessible

- [ ] **Run Database Migration**

  - Execute `migrations/add_soft_delete_to_menu_items.sql`
  - Adds `deleted_at` column for soft delete
  - Updates SKU unique constraint to exclude deleted records

- [ ] **Test Image Upload**

  - Create menu item with image
  - Update menu item with image
  - Verify image URLs work

- [ ] **Test Soft Delete**
  - Delete menu item that's in an order (should work now)
  - Verify deleted items don't appear in GET requests
  - Verify deleted items return 404 when accessed by ID
  - Verify images are kept in MinIO

### Low Priority (Documentation)

- [ ] **Update API Documentation**

  - Update swagger/OpenAPI specs
  - Update README if applicable

- [ ] **Update Team**
  - Share migration guides
  - Conduct team training

---

## Backward Compatibility

### What's Backward Compatible?

✅ **GET Requests**

- All GET endpoints work exactly the same
- Response format unchanged (structure-wise)
- Only `image_url` values changed (now MinIO URLs)

✅ **DELETE Requests**

- Same API signature
- Enhanced functionality (auto-deletes images)

✅ **All Other Resources**

- Categories, Modifiers, Orders, etc. unchanged
- Continue using JSON

### What's NOT Backward Compatible?

❌ **POST /v1/menu-items**

- Must use multipart/form-data
- JSON requests will fail

❌ **PUT /v1/menu-items/:id**

- Must use multipart/form-data
- JSON requests will fail

---

## Recommended Rollout Strategy

### Option 1: Big Bang Migration (Recommended)

1. **Preparation Phase (1-2 days)**

   - Set up MinIO in all environments
   - Update all frontend code
   - Update Postman collections
   - Test thoroughly

2. **Deployment Day**

   - Deploy backend changes
   - Deploy frontend changes
   - Monitor for errors
   - Be ready to rollback

3. **Post-Deployment**
   - Monitor error logs
   - Assist users with issues
   - Update documentation

### Option 2: Gradual Migration (If Needed)

If you need backward compatibility, you could:

1. Create new endpoints:
   - `POST /v1/menu-items/upload` (multipart)
   - `PUT /v1/menu-items/:id/upload` (multipart)
2. Keep old endpoints temporarily:

   - `POST /v1/menu-items` (JSON, deprecated)
   - `PUT /v1/menu-items/:id` (JSON, deprecated)

3. Migrate gradually, then remove old endpoints

**Note:** This requires additional backend work and is NOT currently implemented.

---

## Testing Strategy

### Unit Tests

- [ ] Test multipart parsing
- [ ] Test image upload
- [ ] Test image validation
- [ ] Test error handling

### Integration Tests

- [ ] Test create with image
- [ ] Test create without image
- [ ] Test update with new image
- [ ] Test update without image
- [ ] Test delete (verify image deleted)

### E2E Tests

- [ ] Frontend form submission
- [ ] Image preview
- [ ] Error handling
- [ ] Image URL accessibility

---

## Risk Assessment

### High Risk

- **Breaking changes to active API**
  - Mitigation: Coordinate deployment with frontend
  - Mitigation: Have rollback plan ready

### Medium Risk

- **MinIO dependency**
  - Mitigation: Monitor MinIO health
  - Mitigation: Set up alerts
  - Mitigation: Have backup plan

### Low Risk

- **File size/type validation**
  - Mitigation: Clear error messages
  - Mitigation: Frontend validation

---

## Support Resources

### For Developers

- `MIGRATION_GUIDE.md` - Quick migration steps
- `FRONTEND_IMAGE_UPLOAD_GUIDE.md` - Frontend examples
- `POSTMAN_MINIO_GUIDE.md` - Postman setup

### For DevOps

- `docker-compose.yaml` - MinIO setup
- `configs.example/ENV_VARIABLES.md` - Environment config

### For QA

- Test scenarios in guides
- Postman collection updates

---

## Rollback Plan

If deployment fails:

1. **Immediate Actions**

   - Revert backend deployment
   - Revert frontend deployment
   - Restore database if needed

2. **Data Considerations**

   - Uploaded images remain in MinIO
   - Database `image_url` fields won't break (just URLs)
   - Can manually clean up MinIO later

3. **Communication**
   - Notify team
   - Update status page
   - Post-mortem meeting

---

## Questions?

Contact the development team for assistance with:

- Migration issues
- Integration problems
- MinIO setup
- Testing support
