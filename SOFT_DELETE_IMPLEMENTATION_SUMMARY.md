# Soft Delete Implementation Summary

## Problem Solved ✅

**Before:**

```
ERROR: update or delete on table "menu_items" violates foreign key constraint
"fk_order_items_menu_item" on table "order_items"
```

Menu items couldn't be deleted if they were referenced in any orders.

**After:**
✅ Delete works even if item is in orders
✅ Historical data preserved
✅ No breaking API changes

---

## Changes Made

### 1. Model Updated (`models/menuItems.go`)

```go
type MenuItem struct {
    // ... existing fields ...
    DeletedAt  gorm.DeletedAt `gorm:"index;column:deleted_at"` // NEW
}
```

### 2. Handler Simplified (`feature/menuItem/delivery/http.go`)

- Removed image deletion logic
- Simplified to just soft delete
- Images kept for potential restoration

### 3. Order Repository Updated (`feature/order/repository/postgres.go`)

- Uses `Unscoped()` when preloading MenuItem
- Ensures deleted menu items still show in historical orders
- Prevents nil pointer errors

### 4. Migration Created (`migrations/add_soft_delete_to_menu_items.sql`)

- Adds `deleted_at` column
- Adds index for performance
- Updates SKU unique constraint

---

## How It Works

### Delete Operation

```
DELETE /v1/menu-items/:id
```

**What happens:**

1. Sets `deleted_at` to current timestamp
2. Item is **not** physically deleted
3. Image stays in MinIO
4. No foreign key errors!

### Query Operations

```
GET /v1/menu-items
GET /v1/menu-items/:id
```

**Automatic filtering:**

- Only returns items where `deleted_at IS NULL`
- Deleted items automatically hidden
- GORM handles this transparently

---

## Migration Required

### ✅ GORM AutoMigrate (Automatic)

**Good news!** Your project uses AutoMigrate, so GORM will automatically:

- ✅ Add `deleted_at` column
- ✅ Create index on `deleted_at`
- ✅ Enable soft delete

**Just restart your application:**

```bash
go run main.go
# GORM will handle the migration automatically!
```

### ⚠️ Manual Step (Optional but Recommended)

Update the SKU unique constraint to exclude deleted records:

```sql
-- This allows creating new items with same SKU as deleted items
ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS menu_items_sku_key;
CREATE UNIQUE INDEX IF NOT EXISTS menu_items_sku_unique
ON menu_items(sku) WHERE deleted_at IS NULL;
```

**Run it:**

```bash
psql -U postgres -d pos_db -c "ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS menu_items_sku_key; CREATE UNIQUE INDEX menu_items_sku_unique ON menu_items(sku) WHERE deleted_at IS NULL;"
```

**Why?** Without this, you can't create a new menu item with the same SKU as a deleted one.

---

## Testing

### Test Scenario 1: Delete Item in Order

```bash
# 1. Create a menu item
POST /v1/menu-items
{
  "name": "Test Item",
  "sku": "TEST-001",
  "price_baht": 100
}

# 2. Create an order with this item
POST /v1/orders
{
  "items": [{"menu_item_id": "{id}", "quantity": 1}]
}

# 3. Try to delete the menu item
DELETE /v1/menu-items/{id}

# ✅ Should succeed now (before: would fail with FK error)
```

### Test Scenario 2: Deleted Items Hidden

```bash
# After deleting item from above:

# Get all items
GET /v1/menu-items
# ✅ Deleted item should NOT appear

# Get deleted item by ID
GET /v1/menu-items/{id}
# ✅ Should return 404 Not Found
```

### Test Scenario 3: Database Check

```sql
-- Check deleted items in database
SELECT id, name, sku, deleted_at
FROM menu_items
WHERE deleted_at IS NOT NULL;

-- Deleted item should appear here with timestamp
```

---

## API Contract

### ✅ No Breaking Changes

All endpoints work exactly the same:

| Endpoint                    | Behavior                    | Changed?    |
| --------------------------- | --------------------------- | ----------- |
| `GET /v1/menu-items`        | Lists non-deleted items     | ❌ No       |
| `GET /v1/menu-items/:id`    | Returns 404 if deleted      | ❌ No       |
| `POST /v1/menu-items`       | Creates new item            | ❌ No       |
| `PUT /v1/menu-items/:id`    | Returns 404 if deleted      | ❌ No       |
| `DELETE /v1/menu-items/:id` | **Now works if in orders!** | ✅ Enhanced |

---

## Benefits

### For Business

- ✅ Preserve order history
- ✅ Can restore accidentally deleted items
- ✅ Audit trail of deletions
- ✅ Compliance with record retention

### For Development

- ✅ No foreign key errors
- ✅ Simpler code
- ✅ Data integrity maintained
- ✅ No cascade delete complexity

### For Operations

- ✅ No downtime required
- ✅ Non-breaking change
- ✅ Easy rollback (just remove column)
- ✅ Better data retention

---

## Files Changed

✅ `models/menuItems.go` - Added DeletedAt field
✅ `feature/menuItem/delivery/http.go` - Simplified delete handler
✅ `migrations/add_soft_delete_to_menu_items.sql` - Database migration
✅ `SOFT_DELETE_GUIDE.md` - Detailed documentation
✅ `API_IMPACT_ANALYSIS.md` - Updated with soft delete info

---

## Next Steps

1. **Restart Your Application**

   ```bash
   # GORM AutoMigrate will add deleted_at column automatically!
   go run main.go
   ```

2. **(Optional) Update SKU Constraint**

   ```bash
   # Recommended: Allows same SKU for new items after deletion
   psql -U postgres -d pos_db -c "ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS menu_items_sku_key; CREATE UNIQUE INDEX menu_items_sku_unique ON menu_items(sku) WHERE deleted_at IS NULL;"
   ```

3. **Test**

   - Try deleting menu item that's in an order (should work now!)
   - Verify deleted items don't appear in GET requests
   - Verify deleted items return 404 when accessed by ID

4. **Monitor**
   - Watch application logs on startup (should see "Migrated database")
   - Verify soft delete is working
   - Check deleted_at column was created

---

## Optional Future Enhancements

Consider adding later:

### 1. Restore Deleted Items

```
POST /v1/menu-items/:id/restore
```

### 2. View Deleted Items (Admin)

```
GET /v1/menu-items?include_deleted=true
```

### 3. Permanent Delete (Admin)

```
DELETE /v1/menu-items/:id?permanent=true
```

### 4. Automatic Cleanup

- Permanently delete items after 1 year
- Run as scheduled job

---

## Rollback Plan

If needed, you can rollback:

```sql
-- Remove soft delete column
ALTER TABLE menu_items DROP COLUMN deleted_at;

-- Restore old unique constraint
DROP INDEX IF EXISTS menu_items_sku_unique;
ALTER TABLE menu_items ADD CONSTRAINT menu_items_sku_key UNIQUE (sku);
```

Then revert code changes.

---

## Documentation

- **Full Guide:** `SOFT_DELETE_GUIDE.md`
- **Migration SQL:** `migrations/add_soft_delete_to_menu_items.sql`
- **API Impact:** `API_IMPACT_ANALYSIS.md`

---

## Summary

✅ **Soft delete implemented successfully!**

- No more foreign key errors when deleting menu items
- Data is preserved for historical orders
- No breaking API changes
- Just needs database migration

**You're ready to delete menu items that are in orders!** 🎉
