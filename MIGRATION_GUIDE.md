# Quick Migration Guide: JSON to Multipart Form-Data

## What Changed?

The MenuItem API now accepts **multipart/form-data** instead of **JSON** to support image uploads.

### Before ❌

```json
POST /v1/menu-items
Content-Type: application/json

{
  "category_id": "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
  "name": "ข้าวหมูกรอบพริกเกลือ",
  "sku": "MOU-ASD-PRI-3",
  "price_baht": 70,
  "active": true,
  "image_url": "https://external-url.com/image.jpg"
}
```

### After ✅

```
POST /v1/menu-items
Content-Type: multipart/form-data

Form Fields:
- category_id: 76ea84b7-69da-48bb-bd4d-3a4881f9efcd (Text)
- name: ข้าวหมูกรอบพริกเกลือ (Text)
- sku: MOU-ASD-PRI-3 (Text)
- price_baht: 70 (Text)
- active: true (Text)
- image: [Upload actual image file] (File)
```

---

## Postman Changes

### Step 1: Change Body Type

1. Open your request
2. Go to **Body** tab
3. Change from **raw (JSON)** to **form-data**

### Step 2: Update Fields

Convert your JSON fields to form-data:

| JSON Key        | Form-Data Key | Type     | Value                                  |
| --------------- | ------------- | -------- | -------------------------------------- |
| `category_id`   | `category_id` | Text     | `76ea84b7-69da-48bb-bd4d-3a4881f9efcd` |
| `name`          | `name`        | Text     | `ข้าวหมูกรอบพริกเกลือ`                 |
| `sku`           | `sku`         | Text     | `MOU-ASD-PRI-3`                        |
| `price_baht`    | `price_baht`  | Text     | `70`                                   |
| `active`        | `active`      | Text     | `true`                                 |
| ~~`image_url`~~ | `image`       | **File** | Select image file                      |

### Step 3: Upload Image

1. For the `image` field, click the dropdown
2. Change type from "Text" to "**File**"
3. Click "Select Files" and choose an image
4. Supported formats: JPG, PNG, GIF, WEBP (max 5MB)

### Step 4: Remove Content-Type Header

- Postman will automatically set the correct `Content-Type` with boundary
- **Do not** manually set `Content-Type: multipart/form-data`

### Authentication

- Your cookie-based auth (`token` cookie) is automatically sent
- No need for Authorization header

---

## Frontend Changes

### React Example

#### Before ❌

```typescript
const createMenuItem = async (data) => {
  const response = await fetch("/v1/menu-items", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      name: "ข้าวหมูกรอบพริกเกลือ",
      sku: "MOU-ASD-PRI-3",
      price_baht: 70,
      category_id: "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
      active: true,
      image_url: "https://external-url.com/image.jpg",
    }),
  });
  return response.json();
};
```

#### After ✅

```typescript
const createMenuItem = async (data: {
  name: string;
  sku: string;
  price_baht: number;
  category_id?: string;
  active?: boolean;
  image?: File; // Now accepts File object
}) => {
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

  const response = await fetch("http://localhost:8080/v1/menu-items", {
    method: "POST",
    credentials: "include", // Send cookies
    body: formData,
    // Do NOT set Content-Type header
  });

  return response.json();
};
```

### HTML Form Example

```html
<form id="menuItemForm" enctype="multipart/form-data">
  <input type="text" name="name" required placeholder="ข้าวหมูกรอบพริกเกลือ" />
  <input type="text" name="sku" required placeholder="MOU-ASD-PRI-3" />
  <input type="number" name="price_baht" required placeholder="70" />
  <input
    type="text"
    name="category_id"
    placeholder="76ea84b7-69da-48bb-bd4d-3a4881f9efcd"
  />
  <input type="checkbox" name="active" checked />
  <input type="file" name="image" accept="image/*" />
  <button type="submit">Create</button>
</form>

<script>
  document
    .getElementById("menuItemForm")
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      const formData = new FormData(e.target);

      const response = await fetch("http://localhost:8080/v1/menu-items", {
        method: "POST",
        credentials: "include",
        body: formData,
      });

      const result = await response.json();
      console.log("Created:", result);
      console.log("Image URL:", result.image_url);
    });
</script>
```

### Key Changes Summary

1. ✅ Use `FormData()` instead of `JSON.stringify()`
2. ✅ Add `credentials: 'include'` to send cookies
3. ✅ Remove `Content-Type` header (browser sets it automatically)
4. ✅ Convert `image_url` string → `image` File object
5. ✅ Get `image_url` from **response** (not request)
6. ✅ Convert all values to strings when appending to FormData
7. ✅ Add file input: `<input type="file" name="image" accept="image/*" />`

---

## Response Format (Unchanged)

The response format is the same, but now includes the MinIO URL:

```json
{
  "id": "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
  "name": "ข้าวหมูกรอบพริกเกลือ",
  "price_baht": 70,
  "active": true,
  "image_url": "http://localhost:9000/menu-images/uuid-timestamp.jpg",
  "category": {
    "id": "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
    "name": "อาหารจานเดียว",
    "display_order": 1
  }
}
```

---

## Important Notes

### 1. Image Upload is Optional

You can create menu items **without** images:

```typescript
// No image field - perfectly valid!
formData.append("name", "Item Name");
formData.append("sku", "SKU123");
formData.append("price_baht", "100");
// Don't append image field
```

### 2. Update Without Changing Image

To update a menu item without changing its image:

```typescript
// Simply don't include the image field
formData.append("name", "Updated Name");
formData.append("sku", "SKU123");
formData.append("price_baht", "150");
// No image field = keeps existing image
```

### 3. Cookie Authentication

Your existing cookie-based authentication works automatically:

- No changes needed to auth flow
- Cookies are sent automatically with `credentials: 'include'`
- Make sure users are logged in before creating menu items

### 4. File Validation

The backend validates:

- **File type:** JPG, JPEG, PNG, GIF, WEBP only
- **File size:** Max 5MB

Add frontend validation for better UX:

```typescript
const validateImage = (file: File): string | null => {
  const validTypes = ["image/jpeg", "image/png", "image/gif", "image/webp"];
  if (!validTypes.includes(file.type)) {
    return "Only JPG, PNG, GIF, and WEBP images are allowed";
  }

  if (file.size > 5 * 1024 * 1024) {
    return "Image must be less than 5MB";
  }

  return null;
};
```

---

## Testing Checklist

- [ ] Start MinIO: `docker-compose up -d`
- [ ] Add MinIO environment variables to `configs/.env`
- [ ] Update Postman requests to use form-data
- [ ] Test create with image
- [ ] Test create without image
- [ ] Update frontend to use FormData
- [ ] Test image upload from frontend
- [ ] Verify image URL works (accessible in browser)
- [ ] Test update with new image (old image deleted)
- [ ] Test update without image (existing image preserved)
- [ ] Test delete (image deleted from MinIO)

---

## Need Help?

See detailed guides:

- **Postman:** `POSTMAN_MINIO_GUIDE.md`
- **Frontend:** `FRONTEND_IMAGE_UPLOAD_GUIDE.md`
- **API Reference:** `MINIO_IMAGE_UPLOAD_GUIDE.md`

## Quick Start

1. Add to `configs/.env`:

   ```bash
   MINIO_ENDPOINT=localhost:9000
   MINIO_ROOT_USER=minioadmin
   MINIO_ROOT_PASSWORD=minioadmin
   MINIO_USE_SSL=false
   MINIO_BUCKET_NAME=menu-images
   MINIO_PORT=9000
   MINIO_CONSOLE_PORT=9001
   ```

2. Start services:

   ```bash
   docker-compose up -d
   ```

3. Update your Postman/Frontend as shown above

4. Test!
