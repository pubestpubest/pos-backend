# Frontend Modifier-Category Migration Guide

## Overview

This guide shows how to update your frontend to work with the new modifier-category relationship. Modifiers can now be optionally linked to categories, allowing better organization and filtering.

## What Changed?

### 🔄 Backward Compatible Changes

- ✅ **Existing modifiers still work** - CategoryID is optional and nullable
- ✅ **All existing API calls work** - No breaking changes to request/response structure
- ✅ **New fields are optional** - You can gradually adopt the new features

### 📝 New Features

1. **Category Assignment**: Modifiers can be linked to categories
2. **Category Filtering**: Get modifiers by category
3. **Nested Category Data**: Modifier responses include category information

---

## Updated Data Types

### TypeScript Interfaces

```typescript
// Updated Modifier Response
interface ModifierResponse {
  id: string; // UUID
  category_id?: string | null; // NEW: Optional category UUID
  name: string;
  price_delta_baht: number;
  note: string;
  category?: CategoryResponse; // NEW: Nested category data (if assigned)
}

// Modifier Create/Update Request
interface ModifierRequest {
  name: string; // Required
  category_id?: string | null; // NEW: Optional category UUID
  price_delta_baht?: number; // Optional, defaults to 0
  note?: string; // Optional
}

// Category Response (Updated)
interface CategoryResponse {
  id: string; // UUID
  name: string;
  display_order: number;
  modifiers?: ModifierResponse[]; // NEW: Optional modifiers array
}
```

---

## API Endpoints

### Existing Endpoints (Updated Responses)

#### Get All Modifiers

```http
GET /v1/modifiers
```

**Response:**

```json
[
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
  },
  {
    "id": "223e4567-e89b-12d3-a456-426614174001",
    "category_id": null,
    "name": "No Onions",
    "price_delta_baht": 0,
    "note": "Remove onions",
    "category": null
  }
]
```

#### Get Modifier by ID

```http
GET /v1/modifiers/:id
```

**Response:**

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

#### Create Modifier (Updated)

```http
POST /v1/modifiers
```

**Request Body (NEW - with category):**

```json
{
  "name": "Extra Cheese",
  "category_id": "456e7890-e89b-12d3-a456-426614174111",
  "price_delta_baht": 20,
  "note": "Add extra cheese"
}
```

**Request Body (OLD - still works):**

```json
{
  "name": "Extra Cheese",
  "price_delta_baht": 20,
  "note": "Add extra cheese"
}
```

#### Update Modifier (Updated)

```http
PUT /v1/modifiers/:id
```

**Request Body:**

```json
{
  "name": "Extra Cheese",
  "category_id": "456e7890-e89b-12d3-a456-426614174111",
  "price_delta_baht": 25,
  "note": "Updated note"
}
```

**To remove category assignment:**

```json
{
  "name": "Extra Cheese",
  "category_id": null,
  "price_delta_baht": 25
}
```

### New Endpoints

#### Get Modifiers by Category

```http
GET /v1/categories/:id/modifiers
```

**Response:**

```json
[
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
]
```

---

## React / Next.js Examples

### 1. Fetching Modifiers (Display Category)

```typescript
import { useEffect, useState } from "react";

interface Modifier {
  id: string;
  category_id?: string | null;
  name: string;
  price_delta_baht: number;
  note: string;
  category?: {
    id: string;
    name: string;
    display_order: number;
  };
}

function ModifierList() {
  const [modifiers, setModifiers] = useState<Modifier[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch("http://localhost:8080/v1/modifiers", {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        setModifiers(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching modifiers:", err);
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      <h2>Modifiers</h2>
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Category</th>
            <th>Price Delta (฿)</th>
            <th>Note</th>
          </tr>
        </thead>
        <tbody>
          {modifiers.map((modifier) => (
            <tr key={modifier.id}>
              <td>{modifier.name}</td>
              <td>
                {modifier.category ? (
                  <span className="badge">{modifier.category.name}</span>
                ) : (
                  <span className="text-muted">Uncategorized</span>
                )}
              </td>
              <td>฿{modifier.price_delta_baht}</td>
              <td>{modifier.note || "-"}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
```

### 2. Create Modifier with Category Selection

```typescript
import { useState, useEffect } from "react";

interface Category {
  id: string;
  name: string;
  display_order: number;
}

function CreateModifierForm() {
  const [name, setName] = useState("");
  const [categoryId, setCategoryId] = useState<string>("");
  const [priceDelta, setPriceDelta] = useState(0);
  const [note, setNote] = useState("");
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(false);

  // Fetch categories for the dropdown
  useEffect(() => {
    fetch("http://localhost:8080/v1/categories", {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => setCategories(data))
      .catch((err) => console.error("Error fetching categories:", err));
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    const requestBody: any = {
      name,
      price_delta_baht: priceDelta,
    };

    // Only include category_id if selected
    if (categoryId) {
      requestBody.category_id = categoryId;
    }

    // Only include note if provided
    if (note) {
      requestBody.note = note;
    }

    try {
      const response = await fetch("http://localhost:8080/v1/modifiers", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(requestBody),
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || "Failed to create modifier");
      }

      const result = await response.json();
      console.log("Created modifier:", result);
      alert("Modifier created successfully!");

      // Reset form
      setName("");
      setCategoryId("");
      setPriceDelta(0);
      setNote("");
    } catch (err) {
      console.error("Error creating modifier:", err);
      alert(err instanceof Error ? err.message : "Failed to create modifier");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <h2>Create Modifier</h2>

      <div>
        <label htmlFor="name">Name *</label>
        <input
          id="name"
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>

      <div>
        <label htmlFor="category">Category (Optional)</label>
        <select
          id="category"
          value={categoryId}
          onChange={(e) => setCategoryId(e.target.value)}
        >
          <option value="">-- No Category --</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      <div>
        <label htmlFor="priceDelta">Price Delta (฿)</label>
        <input
          id="priceDelta"
          type="number"
          value={priceDelta}
          onChange={(e) => setPriceDelta(Number(e.target.value))}
        />
      </div>

      <div>
        <label htmlFor="note">Note</label>
        <textarea
          id="note"
          value={note}
          onChange={(e) => setNote(e.target.value)}
        />
      </div>

      <button type="submit" disabled={loading}>
        {loading ? "Creating..." : "Create Modifier"}
      </button>
    </form>
  );
}
```

### 3. Update Modifier (Change Category)

```typescript
async function updateModifier(
  modifierId: string,
  data: {
    name: string;
    category_id?: string | null;
    price_delta_baht?: number;
    note?: string;
  }
) {
  const requestBody: any = {
    name: data.name,
  };

  // Include category_id (can be null to remove category)
  if (data.category_id !== undefined) {
    requestBody.category_id = data.category_id;
  }

  if (data.price_delta_baht !== undefined) {
    requestBody.price_delta_baht = data.price_delta_baht;
  }

  if (data.note !== undefined) {
    requestBody.note = data.note;
  }

  const response = await fetch(
    `http://localhost:8080/v1/modifiers/${modifierId}`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      credentials: "include",
      body: JSON.stringify(requestBody),
    }
  );

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Failed to update modifier");
  }

  return response.json();
}

// Usage examples:

// Update with new category
await updateModifier("modifier-uuid", {
  name: "Extra Cheese",
  category_id: "category-uuid",
  price_delta_baht: 25,
});

// Remove category assignment
await updateModifier("modifier-uuid", {
  name: "Extra Cheese",
  category_id: null,
  price_delta_baht: 25,
});

// Update without changing category
await updateModifier("modifier-uuid", {
  name: "Extra Cheese",
  price_delta_baht: 30,
});
```

### 4. Get Modifiers by Category

```typescript
import { useEffect, useState } from "react";

function CategoryModifiers({ categoryId }: { categoryId: string }) {
  const [modifiers, setModifiers] = useState<Modifier[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetch(`http://localhost:8080/v1/categories/${categoryId}/modifiers`, {
      credentials: "include",
    })
      .then((res) => res.json())
      .then((data) => {
        setModifiers(data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching modifiers:", err);
        setLoading(false);
      });
  }, [categoryId]);

  if (loading) return <div>Loading modifiers...</div>;

  if (modifiers.length === 0) {
    return <div>No modifiers in this category</div>;
  }

  return (
    <div>
      <h3>Modifiers</h3>
      <ul>
        {modifiers.map((modifier) => (
          <li key={modifier.id}>
            {modifier.name} - ฿{modifier.price_delta_baht}
          </li>
        ))}
      </ul>
    </div>
  );
}
```

### 5. Filtering Modifiers by Category (Client-side)

```typescript
import { useState, useMemo } from "react";

function ModifierListWithFilter({ modifiers }: { modifiers: Modifier[] }) {
  const [selectedCategoryId, setSelectedCategoryId] = useState<string>("");

  // Get unique categories from modifiers
  const categories = useMemo(() => {
    const categoryMap = new Map();
    modifiers.forEach((modifier) => {
      if (modifier.category) {
        categoryMap.set(modifier.category.id, modifier.category);
      }
    });
    return Array.from(categoryMap.values());
  }, [modifiers]);

  // Filter modifiers based on selected category
  const filteredModifiers = useMemo(() => {
    if (!selectedCategoryId) return modifiers;

    if (selectedCategoryId === "uncategorized") {
      return modifiers.filter((m) => !m.category_id);
    }

    return modifiers.filter((m) => m.category_id === selectedCategoryId);
  }, [modifiers, selectedCategoryId]);

  return (
    <div>
      <div>
        <label htmlFor="categoryFilter">Filter by Category:</label>
        <select
          id="categoryFilter"
          value={selectedCategoryId}
          onChange={(e) => setSelectedCategoryId(e.target.value)}
        >
          <option value="">All Categories</option>
          <option value="uncategorized">Uncategorized</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      <div>
        <p>
          Showing {filteredModifiers.length} of {modifiers.length} modifiers
        </p>
        <ul>
          {filteredModifiers.map((modifier) => (
            <li key={modifier.id}>
              <strong>{modifier.name}</strong> - ฿{modifier.price_delta_baht}
              {modifier.category && (
                <span className="badge">{modifier.category.name}</span>
              )}
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}
```

---

## Vue.js Example

```typescript
<template>
  <div>
    <h2>Create Modifier</h2>
    <form @submit.prevent="handleSubmit">
      <div>
        <label>Name</label>
        <input v-model="formData.name" required />
      </div>

      <div>
        <label>Category (Optional)</label>
        <select v-model="formData.category_id">
          <option value="">-- No Category --</option>
          <option
            v-for="category in categories"
            :key="category.id"
            :value="category.id"
          >
            {{ category.name }}
          </option>
        </select>
      </div>

      <div>
        <label>Price Delta (฿)</label>
        <input v-model.number="formData.price_delta_baht" type="number" />
      </div>

      <div>
        <label>Note</label>
        <textarea v-model="formData.note"></textarea>
      </div>

      <button type="submit" :disabled="loading">
        {{ loading ? 'Creating...' : 'Create Modifier' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';

interface Category {
  id: string;
  name: string;
  display_order: number;
}

const formData = ref({
  name: '',
  category_id: '',
  price_delta_baht: 0,
  note: ''
});

const categories = ref<Category[]>([]);
const loading = ref(false);

onMounted(async () => {
  try {
    const response = await fetch('http://localhost:8080/v1/categories', {
      credentials: 'include'
    });
    categories.value = await response.json();
  } catch (err) {
    console.error('Error fetching categories:', err);
  }
});

const handleSubmit = async () => {
  loading.value = true;

  const requestBody: any = {
    name: formData.value.name,
    price_delta_baht: formData.value.price_delta_baht
  };

  if (formData.value.category_id) {
    requestBody.category_id = formData.value.category_id;
  }

  if (formData.value.note) {
    requestBody.note = formData.value.note;
  }

  try {
    const response = await fetch('http://localhost:8080/v1/modifiers', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      credentials: 'include',
      body: JSON.stringify(requestBody)
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error);
    }

    const result = await response.json();
    console.log('Created:', result);
    alert('Modifier created successfully!');

    // Reset form
    formData.value = {
      name: '',
      category_id: '',
      price_delta_baht: 0,
      note: ''
    };
  } catch (err) {
    console.error('Error:', err);
    alert(err instanceof Error ? err.message : 'Failed to create modifier');
  } finally {
    loading.value = false;
  }
};
</script>
```

---

## Migration Checklist

### For Existing Applications

- [ ] **Update TypeScript interfaces** to include `category_id` and `category` fields
- [ ] **Update UI** to display category information when available
- [ ] **Handle null categories** gracefully (show "Uncategorized" or similar)
- [ ] **Test existing modifiers** work correctly (they'll have `category_id: null`)

### For New Features

- [ ] **Add category selector** to modifier create/edit forms
- [ ] **Implement category filtering** to show modifiers by category
- [ ] **Use new endpoint** `GET /v1/categories/:id/modifiers` for category-specific views
- [ ] **Add validation** to ensure selected category exists

### Optional Enhancements

- [ ] **Bulk assign categories** to existing modifiers
- [ ] **Show modifier count** per category
- [ ] **Implement drag-and-drop** to change modifier categories
- [ ] **Add category badges** in modifier lists

---

## Important Notes

### 🔒 Authentication Required

All create, update, and delete operations require authentication. Make sure to include `credentials: 'include'` in your fetch requests.

### ⚠️ Validation

- **Category must exist**: If you provide a `category_id`, the backend validates it exists
- **Invalid category**: Returns `400 Bad Request` with error message
- **Name is required**: Modifier name cannot be empty

### 🎯 Best Practices

1. **Always handle null categories**: Not all modifiers will have categories
2. **Use the nested endpoint**: For category-specific pages, use `/v1/categories/:id/modifiers`
3. **Preload categories**: Fetch categories once and cache for dropdown/filter options
4. **Show category in UI**: Display category name in modifier lists for better UX

### 🚀 Performance Tips

1. **Client-side filtering**: For small datasets, filter modifiers client-side
2. **Server-side filtering**: For large datasets, use the category-specific endpoint
3. **Cache category data**: Categories change infrequently, cache them in state management

---

## Testing Your Implementation

### Test Cases

1. ✅ **Create modifier without category** - Should work (backward compatibility)
2. ✅ **Create modifier with category** - Should include category in response
3. ✅ **Update modifier to add category** - Should update successfully
4. ✅ **Update modifier to remove category** (set to null) - Should remove category
5. ✅ **Fetch modifiers** - Should include category data when present
6. ✅ **Fetch by category** - Should return only modifiers in that category
7. ❌ **Create with invalid category ID** - Should return 400 error

### Example Test

```typescript
// Test creating modifier with category
async function testCreateWithCategory() {
  const categoryResponse = await fetch("http://localhost:8080/v1/categories", {
    credentials: "include",
  });
  const categories = await categoryResponse.json();
  const firstCategory = categories[0];

  const response = await fetch("http://localhost:8080/v1/modifiers", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    credentials: "include",
    body: JSON.stringify({
      name: "Test Modifier",
      category_id: firstCategory.id,
      price_delta_baht: 10,
    }),
  });

  const modifier = await response.json();

  console.assert(
    modifier.category_id === firstCategory.id,
    "Category ID should match"
  );
  console.assert(modifier.category !== null, "Category should be populated");
  console.assert(
    modifier.category.id === firstCategory.id,
    "Nested category ID should match"
  );

  console.log("✅ Test passed!");
}
```

---

## Support

If you encounter any issues:

1. Check the backend logs for detailed error messages
2. Verify category IDs are valid UUIDs
3. Ensure authentication cookies are being sent
4. Review this guide for proper request format

For additional help, refer to:

- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md)
- [POSTMAN_TESTING_GUIDE.md](./POSTMAN_TESTING_GUIDE.md)
