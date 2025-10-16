# Modifier-Category Relationship Implementation

## Overview

This document describes the implementation of the modifier-category relationship feature, which allows modifiers to be optionally linked to categories for better organization.

## ✅ Implementation Status: COMPLETE

All changes have been implemented and are **backward compatible**. No database reseed required.

---

## What Was Changed?

### 1. Database Schema (Models)

#### `models/modifiers.go`

- ✅ Added `CategoryID *uuid.UUID` field (nullable - preserves existing data)
- ✅ Added `Category *Category` relationship
- ✅ Configured foreign key with `OnUpdate:SET NULL, OnDelete:SET NULL`

#### `models/categories.go`

- ✅ Added `Modifiers []Modifier` relationship
- ✅ Allows fetching all modifiers for a category

**Migration**: GORM AutoMigrate will automatically add the `category_id` column on next app restart. Existing modifiers will have `category_id = NULL`.

### 2. Request/Response DTOs

#### `request/modifier.go`

- ✅ Added `CategoryID *uuid.UUID` field (optional)
- ✅ Allows creating/updating modifiers with category assignment

#### `response/modifier.go`

- ✅ Added `CategoryID *uuid.UUID` field
- ✅ Added `Category *CategoryResponse` field (nested category data)
- ✅ Uses `omitempty` JSON tags for clean responses

#### `response/category.go`

- ✅ Added `Modifiers []ModifierResponse` field (optional)
- ✅ Allows including modifiers when fetching categories

### 3. Repository Layer

#### `feature/modifier/repository/postgres.go`

- ✅ Added `Preload("Category")` to `GetAllModifiers()`
- ✅ Added `Preload("Category")` to `GetModifierByID()`
- ✅ Added new method `GetModifiersByCategoryID(categoryID uuid.UUID)`

### 4. Use Case Layer

#### `feature/modifier/usecase/usecase.go`

- ✅ Injected `CategoryRepository` for validation
- ✅ Updated constructor to accept both repositories
- ✅ Added category validation in `CreateModifier()`
- ✅ Added category validation in `UpdateModifier()`
- ✅ Implemented `GetModifiersByCategoryID()` method
- ✅ Added helper `mapModifierToResponse()` for consistent mapping
- ✅ Responses now include nested category data when available

### 5. Domain Interfaces

#### `domain/modifier.go`

- ✅ Added `GetModifiersByCategoryID()` to `ModifierRepository` interface
- ✅ Added `GetModifiersByCategoryID()` to `ModifierUsecase` interface

### 6. HTTP Handlers

#### `feature/modifier/delivery/http.go`

- ✅ Added `GetModifiersByCategoryID()` handler
- ✅ Validates category ID parameter
- ✅ Returns 400 for invalid UUIDs

### 7. Routes

#### `routes/modifierRoute.go`

- ✅ Injected `categoryRepository` into modifier usecase
- ✅ Added new public route: `GET /v1/categories/:id/modifiers`
- ✅ All existing routes remain unchanged

---

## API Endpoints

### Existing Endpoints (Updated Responses)

| Method | Endpoint            | Changes                                        |
| ------ | ------------------- | ---------------------------------------------- |
| GET    | `/v1/modifiers`     | Response includes `category_id` and `category` |
| GET    | `/v1/modifiers/:id` | Response includes `category_id` and `category` |
| POST   | `/v1/modifiers`     | Accepts optional `category_id` in request      |
| PUT    | `/v1/modifiers/:id` | Accepts optional `category_id` in request      |
| DELETE | `/v1/modifiers/:id` | No changes                                     |

### New Endpoints

| Method | Endpoint                       | Description                      | Auth Required |
| ------ | ------------------------------ | -------------------------------- | ------------- |
| GET    | `/v1/categories/:id/modifiers` | Get all modifiers for a category | No            |

---

## Data Flow

### Creating a Modifier with Category

```
1. Client sends POST request with category_id
   ↓
2. Handler validates UUID format
   ↓
3. Usecase validates category exists (via CategoryRepository)
   ↓
4. Repository creates modifier with category_id
   ↓
5. Usecase fetches created modifier with Category preloaded
   ↓
6. Response includes nested category data
```

### Fetching Modifiers by Category

```
1. Client sends GET /v1/categories/:id/modifiers
   ↓
2. Handler validates category UUID
   ↓
3. Repository queries modifiers WHERE category_id = :id
   ↓
4. Preloads Category relationship
   ↓
5. Returns array of modifiers with nested category data
```

---

## Database Migration

### Automatic Migration (Recommended)

On next application restart, GORM AutoMigrate will:

1. ✅ Add `category_id` column to `modifiers` table
2. ✅ Set column as nullable UUID
3. ✅ Add foreign key constraint to `categories` table
4. ✅ Configure ON UPDATE SET NULL, ON DELETE SET NULL

**Existing data**: All existing modifiers will have `category_id = NULL` (no data loss).

### Manual Migration (Optional)

If you prefer to run the migration manually:

```sql
-- Add category_id column
ALTER TABLE modifiers
ADD COLUMN category_id UUID NULL;

-- Add foreign key constraint
ALTER TABLE modifiers
ADD CONSTRAINT fk_modifiers_category
FOREIGN KEY (category_id)
REFERENCES categories(id)
ON UPDATE SET NULL
ON DELETE SET NULL;

-- Create index for better query performance
CREATE INDEX idx_modifiers_category_id ON modifiers(category_id);
```

---

## Testing

### 1. Test Backward Compatibility

```bash
# Existing modifiers should still work
curl http://localhost:8080/v1/modifiers

# Response should include category_id: null for old modifiers
```

### 2. Test Creating Modifier with Category

```bash
# Get a category ID first
CATEGORY_ID=$(curl http://localhost:8080/v1/categories | jq -r '.[0].id')

# Create modifier with category
curl -X POST http://localhost:8080/v1/modifiers \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "Extra Cheese",
    "category_id": "'$CATEGORY_ID'",
    "price_delta_baht": 20,
    "note": "Add extra cheese"
  }'
```

### 3. Test Getting Modifiers by Category

```bash
# Get modifiers for a specific category
curl http://localhost:8080/v1/categories/$CATEGORY_ID/modifiers
```

### 4. Test Category Validation

```bash
# Try to create modifier with invalid category ID
curl -X POST http://localhost:8080/v1/modifiers \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "Test",
    "category_id": "00000000-0000-0000-0000-000000000000"
  }'

# Should return 500 with "Category not found" error
```

### 5. Test Update to Remove Category

```bash
# Update modifier to remove category
curl -X PUT http://localhost:8080/v1/modifiers/$MODIFIER_ID \
  -H "Content-Type: application/json" \
  -b cookies.txt \
  -d '{
    "name": "Updated Modifier",
    "category_id": null
  }'
```

---

## Example Responses

### Modifier with Category

```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "category_id": "456e7890-e89b-12d3-a456-426614174111",
  "name": "Extra Cheese",
  "price_delta_baht": 20,
  "note": "Add extra cheese",
  "category": {
    "id": "456e7890-e89b-12d3-a456-426614174111",
    "name": "Toppings",
    "display_order": 1
  }
}
```

### Modifier without Category

```json
{
  "id": "223e4567-e89b-12d3-a456-426614174001",
  "name": "No Onions",
  "price_delta_baht": 0,
  "note": "Remove onions"
}
```

Note: Fields with `null` values are omitted due to `omitempty` JSON tags.

---

## Validation Rules

### Create/Update Modifier

1. ✅ `name` is required
2. ✅ `category_id` is optional
3. ✅ If `category_id` is provided, category must exist
4. ✅ Invalid category UUID returns 400 Bad Request
5. ✅ Non-existent category returns 500 with error message

### Get Modifiers by Category

1. ✅ Category UUID must be valid format
2. ✅ Returns empty array if no modifiers in category
3. ✅ Category doesn't need to exist (just returns empty array)

---

## Performance Considerations

### Database Queries

1. **Get All Modifiers**: 2 queries (modifiers + categories via preload)
2. **Get Modifier by ID**: 2 queries (modifier + category via preload)
3. **Get by Category**: 2 queries (modifiers + category via preload)
4. **Create/Update**: 3 queries (validate category + create/update + fetch with category)

### Optimization Opportunities

1. ✅ Already using `Preload()` for efficient joins
2. ✅ Foreign key index automatically created by GORM
3. 💡 Consider adding composite index on `(category_id, name)` if filtering becomes common
4. 💡 Consider caching category data (changes infrequently)

---

## File Changes Summary

### Modified Files

- ✅ `models/modifiers.go`
- ✅ `models/categories.go`
- ✅ `request/modifier.go`
- ✅ `response/modifier.go`
- ✅ `response/category.go`
- ✅ `feature/modifier/repository/postgres.go`
- ✅ `feature/modifier/usecase/usecase.go`
- ✅ `domain/modifier.go`
- ✅ `feature/modifier/delivery/http.go`
- ✅ `routes/modifierRoute.go`

### New Files

- ✅ `MODIFIER_CATEGORY_IMPLEMENTATION.md` (this file)
- ✅ `FRONTEND_MODIFIER_CATEGORY_MIGRATION_GUIDE.md`

### No Changes Required

- ✅ `database/postgres.go` (AutoMigrate handles schema)
- ✅ Other features remain unaffected
- ✅ Existing seed data works as-is

---

## Next Steps

### For Backend

1. ✅ **Restart application** to trigger GORM AutoMigrate
2. ✅ **Verify migration** succeeded (check logs)
3. ✅ **Test endpoints** using Postman or curl
4. ✅ **Assign categories** to existing modifiers (optional)

### For Frontend

1. 📖 Read `FRONTEND_MODIFIER_CATEGORY_MIGRATION_GUIDE.md`
2. ✅ Update TypeScript interfaces
3. ✅ Update UI to show category information
4. ✅ Add category selector to forms
5. ✅ Handle null categories gracefully
6. ✅ Test all modifier CRUD operations

### For Database Admins

1. ✅ No manual intervention required (AutoMigrate handles it)
2. 💡 Optionally add indexes for better performance
3. 💡 Consider running ANALYZE after migration for query optimization

---

## Rollback Plan

If you need to rollback (unlikely):

### Remove Category Field

```sql
-- Remove foreign key
ALTER TABLE modifiers DROP CONSTRAINT IF EXISTS fk_modifiers_category;

-- Remove column
ALTER TABLE modifiers DROP COLUMN IF EXISTS category_id;
```

### Revert Code

```bash
git revert <commit-hash>
```

---

## Benefits

1. ✅ **Better Organization**: Modifiers can be grouped by category
2. ✅ **Flexible Filtering**: Fetch modifiers by category
3. ✅ **Zero Downtime**: Backward compatible, no breaking changes
4. ✅ **No Data Loss**: Existing modifiers preserved with `category_id = NULL`
5. ✅ **Clean Architecture**: Follows existing patterns (MenuItem-Category)
6. ✅ **Type Safe**: Full TypeScript support for frontend

---

## Architecture Compliance

This implementation follows Clean Architecture principles:

- ✅ **Domain Layer**: Interfaces defined, no external dependencies
- ✅ **Use Case Layer**: Business logic with validation
- ✅ **Repository Layer**: Data access with GORM
- ✅ **Handler Layer**: HTTP handling with Gin
- ✅ **Dependency Rule**: Inner layers don't depend on outer layers
- ✅ **Separation of Concerns**: Each layer has single responsibility

---

## Conclusion

The modifier-category relationship has been successfully implemented with:

- ✅ Zero breaking changes
- ✅ No database reseed required
- ✅ Full backward compatibility
- ✅ Comprehensive frontend guide
- ✅ Following project architecture patterns

Simply restart the application to apply the database migration!
