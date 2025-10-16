# Postman Guide for Menu Item Image Upload

## Authentication

This API uses **cookie-based authentication**. Make sure you have logged in first and the `token` cookie is set.

## Create Menu Item with Image Upload

### Setup

1. **Method:** POST
2. **URL:** `{{base_url}}/v1/menu-items`
3. **Authentication:** Cookie-based (no Authorization header needed)

### Body Configuration

1. Go to **Body** tab
2. Select **form-data** (NOT raw JSON)
3. Add the following fields:

| Key           | Type     | Value                                  | Required |
| ------------- | -------- | -------------------------------------- | -------- |
| `name`        | Text     | `ข้าวหมูกรอบพริกเกลือ`                 | ✅ Yes   |
| `sku`         | Text     | `MOU-ASD-PRI-3`                        | ✅ Yes   |
| `price_baht`  | Text     | `70`                                   | ✅ Yes   |
| `category_id` | Text     | `76ea84b7-69da-48bb-bd4d-3a4881f9efcd` | ❌ No    |
| `active`      | Text     | `true`                                 | ❌ No    |
| `image`       | **File** | Select file from computer              | ❌ No    |

### Important Notes

- **Remove** `Content-Type` header (Postman sets it automatically with boundary)
- **Do NOT** use `image_url` field anymore
- **Upload actual image file** using the `image` field
- The `image` field type must be **File**, not Text

### Screenshots Guide

#### Step 1: Set Method and URL

```
POST {{base_url}}/v1/menu-items
```

#### Step 2: Body Tab

- Click "Body" tab
- Select "form-data" radio button

#### Step 3: Add Form Fields

```
name           | Text | ข้าวหมูกรอบพริกเกลือ
sku            | Text | MOU-ASD-PRI-3
price_baht     | Text | 70
category_id    | Text | 76ea84b7-69da-48bb-bd4d-3a4881f9efcd
active         | Text | true
image          | File | [Select File Button]
```

#### Step 4: Upload Image

- For the `image` row, click the dropdown next to the key name
- Change from "Text" to "File"
- Click "Select Files" button that appears
- Choose an image file (JPG, PNG, GIF, WEBP - max 5MB)

## Update Menu Item

### Setup

1. **Method:** PUT
2. **URL:** `{{base_url}}/v1/menu-items/{{menu_item_id}}`
3. **Authentication:** Cookie-based

### Body Configuration

Same as Create (use form-data with same fields)

### Behavior

- If you include `image` field with a file: uploads new image, deletes old one
- If you **don't** include `image` field: keeps existing image
- To update without changing image: just omit the `image` field

## Example Requests

### Create with Image

```
POST http://localhost:8080/v1/menu-items

Form Data:
- name: ข้าวหมูกรอบพริกเกลือ
- sku: MOU-ASD-PRI-3
- price_baht: 70
- category_id: 76ea84b7-69da-48bb-bd4d-3a4881f9efcd
- active: true
- image: food-image.jpg (file upload)

Cookie: token=your-session-token
```

### Create without Image

```
POST http://localhost:8080/v1/menu-items

Form Data:
- name: ข้าวหมูกรอบพริกเกลือ
- sku: MOU-ASD-PRI-3
- price_baht: 70
- category_id: 76ea84b7-69da-48bb-bd4d-3a4881f9efcd
- active: true
(no image field)

Cookie: token=your-session-token
```

### Update with New Image

```
PUT http://localhost:8080/v1/menu-items/76ea84b7-69da-48bb-bd4d-3a4881f9efcd

Form Data:
- name: ข้าวหมูกรอบพริกเกลือพิเศษ
- sku: MOU-ASD-PRI-3
- price_baht: 80
- category_id: 76ea84b7-69da-48bb-bd4d-3a4881f9efcd
- active: true
- image: new-food-image.jpg (file upload)

Cookie: token=your-session-token
```

### Update without Changing Image

```
PUT http://localhost:8080/v1/menu-items/76ea84b7-69da-48bb-bd4d-3a4881f9efcd

Form Data:
- name: ข้าวหมูกรอบพริกเกลือพิเศษ
- sku: MOU-ASD-PRI-3
- price_baht: 80
- category_id: 76ea84b7-69da-48bb-bd4d-3a4881f9efcd
- active: true
(no image field - keeps existing image)

Cookie: token=your-session-token
```

## Response Format

### Success Response (201 Created / 200 OK)

```json
{
  "id": "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
  "name": "ข้าวหมูกรอบพริกเกลือ",
  "price_baht": 70,
  "active": true,
  "image_url": "http://localhost:9000/menu-images/550e8400-1234-5678-9012-1697486400.jpg",
  "category": {
    "id": "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
    "name": "อาหารจานเดียว",
    "display_order": 1
  }
}
```

### Error Responses

#### Missing Required Fields (400)

```json
{
  "error": "Name, SKU, and price_baht are required"
}
```

#### Invalid File Type (500)

```json
{
  "error": "Invalid file type. Only image files are allowed (jpg, jpeg, png, gif, webp)"
}
```

#### File Too Large (500)

```json
{
  "error": "File size exceeds maximum allowed size of 5MB"
}
```

#### Unauthorized (401)

```json
{
  "error": "Authentication cookie required"
}
```

## Testing Checklist

- [ ] Login first to get session cookie
- [ ] Verify cookie is included in request
- [ ] Change from JSON to form-data
- [ ] Use Text type for all fields except image
- [ ] Use File type for image field
- [ ] Test with valid image file (< 5MB)
- [ ] Test without image (should work)
- [ ] Test update with new image
- [ ] Test update without image field (should preserve)
- [ ] Verify image URL in response
- [ ] Access image URL in browser (should work)

## Common Issues

### Issue: "Invalid request body"

**Solution:** Make sure you're using **form-data**, not JSON

### Issue: "Authentication cookie required"

**Solution:** Login first to get the `token` cookie

### Issue: Image not uploading

**Solution:**

1. Change field type from Text to File
2. Check file size (< 5MB)
3. Check file format (JPG, PNG, GIF, WEBP)

### Issue: Getting old JSON error

**Solution:** Clear Postman cache and restart

## Migration from Old JSON API

### Old Way (Before)

```json
POST /v1/menu-items
Content-Type: application/json

{
  "name": "ข้าวหมูกรอบพริกเกลือ",
  "sku": "MOU-ASD-PRI-3",
  "price_baht": 70,
  "category_id": "76ea84b7-69da-48bb-bd4d-3a4881f9efcd",
  "active": true,
  "image_url": "https://external-url.com/image.jpg"
}
```

### New Way (After)

```
POST /v1/menu-items
Content-Type: multipart/form-data

Form Fields:
- name: ข้าวหมูกรอบพริกเกลือ
- sku: MOU-ASD-PRI-3
- price_baht: 70
- category_id: 76ea84b7-69da-48bb-bd4d-3a4881f9efcd
- active: true
- image: [Upload actual file]
```

**Key Difference:**

- ❌ No more `image_url` field with external URL
- ✅ Upload actual image file, get MinIO URL back in response
