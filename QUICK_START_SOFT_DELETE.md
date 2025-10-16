# Quick Start: Soft Delete for Menu Items

## TL;DR

✅ **Just restart your app!** GORM AutoMigrate handles most of it automatically.

```bash
go run main.go
```

That's it! Soft delete is now active. 🎉

---

## What Happens Automatically

When you restart your application, GORM AutoMigrate will:

✅ Add `deleted_at` column to `menu_items` table  
✅ Create index on `deleted_at`  
✅ Enable soft delete functionality

**You'll see this in your logs:**

```
[database]: Connected to database
[database]: Migrated database
```

---

## Optional Manual Step

**Update SKU unique constraint** (recommended but not required):

```bash
psql -U postgres -d pos_db -c "ALTER TABLE menu_items DROP CONSTRAINT IF EXISTS menu_items_sku_key; CREATE UNIQUE INDEX menu_items_sku_unique ON menu_items(sku) WHERE deleted_at IS NULL;"
```

**Why?** This allows you to create a new menu item with the same SKU as a deleted one.

**If you skip this:** You can still use soft delete, but you can't reuse SKUs from deleted items.

---

## How to Test

### 1. Start your application

```bash
go run main.go
```

### 2. Try to delete a menu item

```bash
DELETE /v1/menu-items/{id}
```

✅ Should work even if the item is in orders!

### 3. Verify it's hidden

```bash
GET /v1/menu-items/{id}
```

✅ Should return 404 (item is soft-deleted)

### 4. Check database

```sql
SELECT id, name, sku, deleted_at FROM menu_items WHERE deleted_at IS NOT NULL;
```

✅ You should see your deleted item with a timestamp

---

## What Changed

### Before (Hard Delete)

```
DELETE /v1/menu-items/{id}
❌ ERROR: violates foreign key constraint
```

### After (Soft Delete)

```
DELETE /v1/menu-items/{id}
✅ 200 OK
{
  "message": "Menu item deleted successfully"
}
```

---

## Summary

| What                    | Automatic? | Action             |
| ----------------------- | ---------- | ------------------ |
| Add `deleted_at` column | ✅ Yes     | Just restart app   |
| Add index               | ✅ Yes     | Just restart app   |
| Enable soft delete      | ✅ Yes     | Just restart app   |
| Update SKU constraint   | ⚠️ No      | Run SQL (optional) |

**Bottom line: Just restart your app and it works!** 🚀

---

## Full Documentation

For more details, see:

- `SOFT_DELETE_IMPLEMENTATION_SUMMARY.md` - Overview
- `SOFT_DELETE_GUIDE.md` - Complete guide
- `API_IMPACT_ANALYSIS.md` - API changes
