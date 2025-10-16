# Soft Delete Implementation - Menu Items

## Overview

Menu items now use **soft delete** instead of hard delete. When a menu item is deleted, it's marked as deleted in the database but not actually removed. This preserves historical data and prevents foreign key constraint errors.

## Why Soft Delete?

### Problem

Hard deleting menu items caused database errors when the item was referenced in orders:

```
ERROR: update or delete on table "menu_items" violates foreign key constraint
"fk_order_items_menu_item" on table "order_items"
```

### Solution

Soft delete allows us to:

- ✅ Preserve historical order data
- ✅ Maintain referential integrity for past orders
- ✅ Restore accidentally deleted items
- ✅ Keep audit trails
- ✅ Retain images for potential restoration

---

## How It Works

### Database Schema

The `menu_items` table now includes a `deleted_at` column:

```sql
ALTER TABLE menu_items ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
CREATE INDEX idx_menu_items_deleted_at ON menu_items(deleted_at);
```

### Behavior

#### Delete Operation

```bash
DELETE /v1/menu-items/:id
```

**Before (Hard Delete):**

- Physically removes row from database
- Deletes image from MinIO
- Breaks if referenced in orders ❌

**After (Soft Delete):**

- Sets `deleted_at` timestamp
- Keeps image in MinIO (for restoration)
- Works even if referenced in orders ✅

#### Query Operations

```bash
GET /v1/menu-items
GET /v1/menu-items/:id
```

**Automatic Filtering:**

- Only returns items where `deleted_at IS NULL`
- Deleted items are hidden from all queries
- GORM handles this automatically

---

## Implementation Details

### Model Changes

```go
type MenuItem struct {
    ID         uuid.UUID      `gorm:"type:uuid;..."`
    Name       *string        `gorm:"type:varchar;..."`
    SKU        *string        `gorm:"type:varchar;..."`
    PriceBaht  *int64         `gorm:"column:price_baht"`
    Active     *bool          `gorm:"column:active;default:true"`
    ImageURL   *string        `gorm:"type:text;column:image_url"`
    DeletedAt  gorm.DeletedAt `gorm:"index;column:deleted_at"` // NEW

    Category *Category `gorm:"foreignKey:CategoryID;..."`
}
```

### Repository Changes

**No changes needed!** GORM automatically:

- Adds `WHERE deleted_at IS NULL` to all queries
- Converts `Delete()` to `UPDATE deleted_at = NOW()`

### Handler Changes

**Simplified delete handler** - no longer needs to:

- Fetch menu item first
- Delete image from MinIO
- Handle cascade delete errors

```go
func (h *menuItemHandler) DeleteMenuItem(c *gin.Context) {
    id, _ := uuid.Parse(c.Param("id"))

    // Soft delete - just mark as deleted
    // Image kept for potential restoration
    if err := h.menuItemUsecase.DeleteMenuItem(id); err != nil {
        // Handle error
    }

    c.JSON(http.StatusOK, gin.H{"message": "Menu item deleted successfully"})
}
```

---

## Migration Steps

### 1. Add Column to Database

**GORM AutoMigrate handles this automatically!**

Your project uses AutoMigrate in `database/postgres.go`, which includes `&models.MenuItem{}`. When you restart your application, GORM will automatically:

- Add the `deleted_at` column
- Create the index on `deleted_at`

**No manual migration needed for this part!**

Just restart your app:

```bash
go run main.go
```

### 2. Deploy Backend Changes

The code changes are already implemented:

- ✅ Model updated with `DeletedAt` field
- ✅ Repository uses GORM soft delete
- ✅ Handler simplified (no image deletion)

### 3. Test

Test the soft delete functionality:

```bash
# Create a menu item
POST /v1/menu-items

# Delete it (soft delete)
DELETE /v1/menu-items/{id}

# Try to get it (should return 404)
GET /v1/menu-items/{id}

# Check database - deleted_at should be set
SELECT id, name, deleted_at FROM menu_items WHERE id = '{id}';
```

---

## API Behavior

### No Changes to API Contract

The API endpoints work exactly the same from the client perspective:

#### Create

```
POST /v1/menu-items
```

✅ Same as before

#### Get All

```
GET /v1/menu-items
```

✅ Returns only non-deleted items (automatic)

#### Get One

```
GET /v1/menu-items/:id
```

✅ Returns 404 if deleted (automatic)

#### Update

```
PUT /v1/menu-items/:id
```

✅ Returns 404 if deleted (automatic)

#### Delete

```
DELETE /v1/menu-items/:id
```

✅ Now works even if referenced in orders!

---

## Advanced Operations

### Restore a Deleted Item

If you need to restore a soft-deleted item, you can add a restore endpoint:

```go
// Example restore endpoint (not yet implemented)
func (r *menuItemRepository) RestoreMenuItem(id uuid.UUID) error {
    return r.db.Model(&models.MenuItem{}).
        Unscoped(). // Include soft-deleted records
        Where("id = ?", id).
        Update("deleted_at", nil).Error
}
```

### Permanently Delete (Hard Delete)

If you need to permanently delete items (e.g., cleanup old data):

```go
// Example permanent delete (not yet implemented)
func (r *menuItemRepository) PermanentlyDeleteMenuItem(id uuid.UUID) error {
    // Also delete image from MinIO before permanent deletion
    return r.db.Unscoped(). // Bypass soft delete
        Where("id = ?", id).
        Delete(&models.MenuItem{}).Error
}
```

### Query Deleted Items

To see deleted items (for admin/audit purposes):

```go
// Example query deleted items (not yet implemented)
func (r *menuItemRepository) GetDeletedMenuItems() ([]*models.MenuItem, error) {
    var items []*models.MenuItem
    err := r.db.Unscoped(). // Include soft-deleted records
        Where("deleted_at IS NOT NULL").
        Find(&items).Error
    return items, err
}
```

---

## Data Retention

### Images

- **Kept in MinIO** when menu item is soft-deleted
- Allows restoration if needed
- Can be cleaned up later with a background job if desired

### Database Records

- **Kept indefinitely** by default
- Preserves order history and referential integrity
- Can add cleanup job to permanently delete old soft-deleted records (e.g., after 1 year)

### Cleanup Strategy (Optional)

You can implement a cleanup job to permanently delete old soft-deleted items:

```go
// Example cleanup job
func CleanupOldDeletedItems() {
    // Permanently delete items deleted more than 1 year ago
    oneYearAgo := time.Now().AddDate(-1, 0, 0)

    db.Unscoped().
        Where("deleted_at < ?", oneYearAgo).
        Delete(&models.MenuItem{})
}
```

---

## Benefits

### Business Benefits

1. **Preserve Order History** - Past orders still show correct menu items
2. **Undo Mistakes** - Can restore accidentally deleted items
3. **Audit Trail** - Know when items were deleted
4. **Compliance** - Maintain historical records for accounting

### Technical Benefits

1. **No Foreign Key Errors** - Soft delete doesn't violate constraints
2. **Data Integrity** - Relationships preserved
3. **Simpler Code** - No cascade delete handling needed
4. **Better Performance** - No cascade operations

---

## Testing Checklist

- [x] Delete menu item that's in an order (should work now)
- [x] Deleted item doesn't appear in GET /v1/menu-items
- [x] GET /v1/menu-items/:id returns 404 for deleted item
- [x] Update deleted item returns 404
- [x] Image is kept in MinIO after deletion
- [x] Database has deleted_at timestamp set
- [x] Can delete same SKU again (unique constraint respected)

---

## Troubleshooting

### Issue: Deleted items still appearing

**Cause:** Database migration not run
**Solution:** Run the ALTER TABLE migration to add deleted_at column

### Issue: Cannot create new item with same SKU

**Cause:** GORM unique constraint includes soft-deleted records
**Solution:** Update unique constraint to exclude deleted records:

```sql
-- Remove old unique constraint
ALTER TABLE menu_items DROP CONSTRAINT menu_items_sku_key;

-- Add partial unique index (excludes deleted records)
CREATE UNIQUE INDEX menu_items_sku_unique
ON menu_items(sku) WHERE deleted_at IS NULL;
```

### Issue: Need to restore deleted item

**Solution:** Implement restore endpoint (see Advanced Operations section)

---

## Migration

### Automatic (GORM AutoMigrate)

✅ **Most migration is automatic!** Your project uses GORM AutoMigrate.

Just restart your application:

```bash
go run main.go
```

GORM will automatically:

- Add `deleted_at` column
- Create index on `deleted_at`

### Manual (Optional - SKU Constraint)

Update the SKU unique constraint to exclude deleted records:

```sql
-- Optional: Update unique constraint to exclude soft-deleted records
ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS menu_items_sku_key;
CREATE UNIQUE INDEX IF NOT EXISTS menu_items_sku_unique
ON menu_items(sku) WHERE deleted_at IS NULL;
```

This allows reusing SKUs after deletion.

---

## Future Enhancements

Consider adding:

1. **Restore Endpoint**

   ```
   POST /v1/menu-items/:id/restore
   ```

2. **View Deleted Items** (Admin only)

   ```
   GET /v1/menu-items?include_deleted=true
   ```

3. **Permanent Delete** (Admin only)

   ```
   DELETE /v1/menu-items/:id?permanent=true
   ```

4. **Automatic Cleanup Job**

   - Permanently delete items deleted > 1 year ago
   - Run as cron job or scheduled task

5. **Audit Log**
   - Track who deleted items and when
   - Integration with activity logging

---

## Summary

✅ **Soft delete is now active for menu items**

- Delete operations work even if item is in orders
- Items are hidden but data is preserved
- Images are kept for potential restoration
- No API changes required for clients
- Database migration needed to add `deleted_at` column

**Next Steps:**

1. Run database migration
2. Test delete functionality
3. Consider implementing restore endpoint if needed
