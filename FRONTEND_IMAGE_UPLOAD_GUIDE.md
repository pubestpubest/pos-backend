# Frontend Image Upload Guide

## Overview

This guide shows how to update your frontend to work with the new multipart/form-data API for menu item image uploads.

## Authentication

The API uses **cookie-based authentication**. The session token cookie is automatically sent with requests if you use `credentials: 'include'`.

---

## React / Next.js Example

### Create Menu Item with Image

```typescript
import { useState } from "react";

interface CreateMenuItemData {
  name: string;
  sku: string;
  price_baht: number;
  category_id?: string;
  active?: boolean;
  image?: File;
}

async function createMenuItem(data: CreateMenuItemData) {
  // Create FormData object
  const formData = new FormData();

  // Add required fields
  formData.append("name", data.name);
  formData.append("sku", data.sku);
  formData.append("price_baht", data.price_baht.toString());

  // Add optional fields
  if (data.category_id) {
    formData.append("category_id", data.category_id);
  }

  if (data.active !== undefined) {
    formData.append("active", data.active.toString());
  }

  // Add image file if provided
  if (data.image) {
    formData.append("image", data.image);
  }

  // Send request
  const response = await fetch("http://localhost:8080/v1/menu-items", {
    method: "POST",
    credentials: "include", // Important: sends cookies
    body: formData,
    // Do NOT set Content-Type header - browser sets it automatically with boundary
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Failed to create menu item");
  }

  return response.json();
}

// Usage in React component
function CreateMenuItemForm() {
  const [name, setName] = useState("");
  const [sku, setSku] = useState("");
  const [price, setPrice] = useState(0);
  const [categoryId, setCategoryId] = useState("");
  const [active, setActive] = useState(true);
  const [imageFile, setImageFile] = useState<File | null>(null);
  const [loading, setLoading] = useState(false);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0]);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const result = await createMenuItem({
        name,
        sku,
        price_baht: price,
        category_id: categoryId,
        active,
        image: imageFile || undefined,
      });

      console.log("Created:", result);
      alert("Menu item created successfully!");

      // Display the uploaded image
      if (result.image_url) {
        console.log("Image URL:", result.image_url);
      }
    } catch (error) {
      console.error("Error:", error);
      alert(error.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label>Name:</label>
        <input
          type="text"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>

      <div>
        <label>SKU:</label>
        <input
          type="text"
          value={sku}
          onChange={(e) => setSku(e.target.value)}
          required
        />
      </div>

      <div>
        <label>Price (Baht):</label>
        <input
          type="number"
          value={price}
          onChange={(e) => setPrice(Number(e.target.value))}
          required
        />
      </div>

      <div>
        <label>Category ID:</label>
        <input
          type="text"
          value={categoryId}
          onChange={(e) => setCategoryId(e.target.value)}
        />
      </div>

      <div>
        <label>Active:</label>
        <input
          type="checkbox"
          checked={active}
          onChange={(e) => setActive(e.target.checked)}
        />
      </div>

      <div>
        <label>Image:</label>
        <input
          type="file"
          accept="image/jpeg,image/png,image/gif,image/webp"
          onChange={handleImageChange}
        />
        {imageFile && <p>Selected: {imageFile.name}</p>}
      </div>

      <button type="submit" disabled={loading}>
        {loading ? "Creating..." : "Create Menu Item"}
      </button>
    </form>
  );
}
```

### Update Menu Item

```typescript
async function updateMenuItem(id: string, data: CreateMenuItemData) {
  const formData = new FormData();

  formData.append("name", data.name);
  formData.append("sku", data.sku);
  formData.append("price_baht", data.price_baht.toString());

  if (data.category_id) {
    formData.append("category_id", data.category_id);
  }

  if (data.active !== undefined) {
    formData.append("active", data.active.toString());
  }

  // Only add image if user selected a new one
  if (data.image) {
    formData.append("image", data.image);
  }
  // If no image provided, existing image is preserved

  const response = await fetch(`http://localhost:8080/v1/menu-items/${id}`, {
    method: "PUT",
    credentials: "include",
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Failed to update menu item");
  }

  return response.json();
}
```

---

## Vue.js Example

### Composition API

```vue
<template>
  <form @submit.prevent="handleSubmit">
    <div>
      <label>Name:</label>
      <input v-model="formData.name" type="text" required />
    </div>

    <div>
      <label>SKU:</label>
      <input v-model="formData.sku" type="text" required />
    </div>

    <div>
      <label>Price (Baht):</label>
      <input v-model.number="formData.price_baht" type="number" required />
    </div>

    <div>
      <label>Category ID:</label>
      <input v-model="formData.category_id" type="text" />
    </div>

    <div>
      <label>Active:</label>
      <input v-model="formData.active" type="checkbox" />
    </div>

    <div>
      <label>Image:</label>
      <input
        type="file"
        accept="image/jpeg,image/png,image/gif,image/webp"
        @change="handleFileChange"
      />
      <p v-if="selectedFile">Selected: {{ selectedFile.name }}</p>
    </div>

    <button type="submit" :disabled="loading">
      {{ loading ? "Creating..." : "Create Menu Item" }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref } from "vue";

const formData = ref({
  name: "",
  sku: "",
  price_baht: 0,
  category_id: "",
  active: true,
});

const selectedFile = ref<File | null>(null);
const loading = ref(false);

const handleFileChange = (event: Event) => {
  const target = event.target as HTMLInputElement;
  if (target.files && target.files[0]) {
    selectedFile.value = target.files[0];
  }
};

const handleSubmit = async () => {
  loading.value = true;

  try {
    const formDataObj = new FormData();

    formDataObj.append("name", formData.value.name);
    formDataObj.append("sku", formData.value.sku);
    formDataObj.append("price_baht", formData.value.price_baht.toString());

    if (formData.value.category_id) {
      formDataObj.append("category_id", formData.value.category_id);
    }

    formDataObj.append("active", formData.value.active.toString());

    if (selectedFile.value) {
      formDataObj.append("image", selectedFile.value);
    }

    const response = await fetch("http://localhost:8080/v1/menu-items", {
      method: "POST",
      credentials: "include",
      body: formDataObj,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "Failed to create menu item");
    }

    const result = await response.json();
    console.log("Created:", result);
    alert("Menu item created successfully!");

    // Reset form
    formData.value = {
      name: "",
      sku: "",
      price_baht: 0,
      category_id: "",
      active: true,
    };
    selectedFile.value = null;
  } catch (error) {
    console.error("Error:", error);
    alert(error.message);
  } finally {
    loading.value = false;
  }
};
</script>
```

---

## Vanilla JavaScript / HTML Example

```html
<!DOCTYPE html>
<html>
  <head>
    <title>Create Menu Item</title>
  </head>
  <body>
    <form id="menuItemForm">
      <div>
        <label>Name:</label>
        <input type="text" id="name" name="name" required />
      </div>

      <div>
        <label>SKU:</label>
        <input type="text" id="sku" name="sku" required />
      </div>

      <div>
        <label>Price (Baht):</label>
        <input type="number" id="price_baht" name="price_baht" required />
      </div>

      <div>
        <label>Category ID:</label>
        <input type="text" id="category_id" name="category_id" />
      </div>

      <div>
        <label>Active:</label>
        <input type="checkbox" id="active" name="active" checked />
      </div>

      <div>
        <label>Image:</label>
        <input
          type="file"
          id="image"
          name="image"
          accept="image/jpeg,image/png,image/gif,image/webp"
        />
      </div>

      <button type="submit">Create Menu Item</button>
    </form>

    <script>
      document
        .getElementById("menuItemForm")
        .addEventListener("submit", async (e) => {
          e.preventDefault();

          const formData = new FormData();

          // Get form values
          const name = document.getElementById("name").value;
          const sku = document.getElementById("sku").value;
          const price_baht = document.getElementById("price_baht").value;
          const category_id = document.getElementById("category_id").value;
          const active = document.getElementById("active").checked;
          const imageFile = document.getElementById("image").files[0];

          // Append to FormData
          formData.append("name", name);
          formData.append("sku", sku);
          formData.append("price_baht", price_baht);

          if (category_id) {
            formData.append("category_id", category_id);
          }

          formData.append("active", active.toString());

          if (imageFile) {
            formData.append("image", imageFile);
          }

          try {
            const response = await fetch(
              "http://localhost:8080/v1/menu-items",
              {
                method: "POST",
                credentials: "include", // Sends cookies
                body: formData,
                // Do NOT set Content-Type header
              }
            );

            if (!response.ok) {
              const error = await response.json();
              throw new Error(error.error || "Failed to create menu item");
            }

            const result = await response.json();
            console.log("Created:", result);
            alert("Menu item created successfully!");

            // Display image
            if (result.image_url) {
              console.log("Image URL:", result.image_url);
              // You can display the image:
              // const img = document.createElement('img');
              // img.src = result.image_url;
              // document.body.appendChild(img);
            }

            // Reset form
            e.target.reset();
          } catch (error) {
            console.error("Error:", error);
            alert(error.message);
          }
        });
    </script>
  </body>
</html>
```

---

## Axios Example

```typescript
import axios from "axios";

// Configure axios to send cookies
const api = axios.create({
  baseURL: "http://localhost:8080/v1",
  withCredentials: true, // Important: sends cookies
});

async function createMenuItem(data: {
  name: string;
  sku: string;
  price_baht: number;
  category_id?: string;
  active?: boolean;
  image?: File;
}) {
  const formData = new FormData();

  formData.append("name", data.name);
  formData.append("sku", data.sku);
  formData.append("price_baht", data.price_baht.toString());

  if (data.category_id) {
    formData.append("category_id", data.category_id);
  }

  if (data.active !== undefined) {
    formData.append("active", data.active.toString());
  }

  if (data.image) {
    formData.append("image", data.image);
  }

  const response = await api.post("/menu-items", formData, {
    headers: {
      // Axios sets Content-Type automatically for FormData
      // 'Content-Type': 'multipart/form-data' // Not needed
    },
  });

  return response.data;
}

async function updateMenuItem(
  id: string,
  data: {
    name: string;
    sku: string;
    price_baht: number;
    category_id?: string;
    active?: boolean;
    image?: File;
  }
) {
  const formData = new FormData();

  formData.append("name", data.name);
  formData.append("sku", data.sku);
  formData.append("price_baht", data.price_baht.toString());

  if (data.category_id) {
    formData.append("category_id", data.category_id);
  }

  if (data.active !== undefined) {
    formData.append("active", data.active.toString());
  }

  if (data.image) {
    formData.append("image", data.image);
  }

  const response = await api.put(`/menu-items/${id}`, formData);
  return response.data;
}

// Usage
const fileInput = document.querySelector(
  'input[type="file"]'
) as HTMLInputElement;
const file = fileInput.files?.[0];

createMenuItem({
  name: "ข้าวหมูกรอบพริกเกลือ",
  sku: "MOU-ASD-PRI-3",
  price_baht: 70,
  category_id: "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
  active: true,
  image: file,
})
  .then((result) => {
    console.log("Success:", result);
    console.log("Image URL:", result.image_url);
  })
  .catch((error) => {
    console.error("Error:", error);
  });
```

---

## Important Notes

### 1. Cookie Authentication

```typescript
// Always include credentials to send cookies
fetch(url, {
  credentials: "include", // fetch
});

// For axios
axios.create({
  withCredentials: true, // axios
});
```

### 2. Do NOT Set Content-Type Header

```typescript
// ❌ WRONG - Do not set Content-Type manually
fetch(url, {
  headers: {
    "Content-Type": "multipart/form-data", // Wrong!
  },
  body: formData,
});

// ✅ CORRECT - Let browser set it automatically
fetch(url, {
  // No Content-Type header
  body: formData,
});
```

The browser automatically sets `Content-Type: multipart/form-data; boundary=...` with the correct boundary string.

### 3. Optional Image Upload

```typescript
// Image is optional - you can create menu items without images
const formData = new FormData();
formData.append("name", "Item Name");
formData.append("sku", "SKU123");
formData.append("price_baht", "100");
// No image field - perfectly valid!
```

### 4. Updating Without Changing Image

```typescript
// To update menu item without changing image, simply don't include the image field
const formData = new FormData();
formData.append("name", "Updated Name");
formData.append("sku", "SKU123");
formData.append("price_baht", "150");
// No image field - existing image is preserved
```

### 5. File Validation

```typescript
const validateImage = (file: File): string | null => {
  // Check file type
  const validTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];
  if (!validTypes.includes(file.type)) {
    return "Invalid file type. Only JPG, PNG, GIF, and WEBP are allowed.";
  }

  // Check file size (5MB)
  const maxSize = 5 * 1024 * 1024;
  if (file.size > maxSize) {
    return "File size must be less than 5MB.";
  }

  return null; // Valid
};

// Usage
const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
  const file = e.target.files?.[0];
  if (file) {
    const error = validateImage(file);
    if (error) {
      alert(error);
      e.target.value = ""; // Clear input
      return;
    }
    setImageFile(file);
  }
};
```

---

## Migration Checklist

- [ ] Update API calls from JSON to FormData
- [ ] Change `Content-Type` from `application/json` to auto (remove header)
- [ ] Add `credentials: 'include'` or `withCredentials: true`
- [ ] Remove `image_url` field from requests
- [ ] Add file input for image upload
- [ ] Handle `image_url` from response to display images
- [ ] Add image validation (type, size)
- [ ] Test create with image
- [ ] Test create without image
- [ ] Test update with new image
- [ ] Test update without changing image
- [ ] Handle error responses
- [ ] Update TypeScript interfaces if using TypeScript

---

## Common Issues

### Issue: "Authentication cookie required"

**Cause:** Cookies not being sent with request
**Solution:** Add `credentials: 'include'` (fetch) or `withCredentials: true` (axios)

### Issue: "Failed to parse form data"

**Cause:** Sending JSON instead of FormData
**Solution:** Use `new FormData()` instead of `JSON.stringify()`

### Issue: "Invalid file type"

**Cause:** File is not an image or unsupported format
**Solution:** Validate file type before upload

### Issue: CORS error

**Cause:** Backend not configured for credentials
**Solution:** Backend should have CORS configured (already done in your project)

### Issue: Image uploaded but URL is empty

**Cause:** Check backend logs for MinIO errors
**Solution:** Ensure MinIO is running and environment variables are set
