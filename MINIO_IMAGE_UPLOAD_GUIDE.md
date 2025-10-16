# MinIO Image Upload Guide

This guide explains how to use the MinIO object storage integration for uploading menu item images.

## Overview

The application now supports image uploads for menu items using MinIO, an S3-compatible object storage service. Images are stored in MinIO and their URLs are saved in the database.

## Setup

### 1. Environment Configuration

Add the following environment variables to your `configs/.env` file:

```bash
# MinIO Configuration
MINIO_ENDPOINT=localhost:9000
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=menu-images
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
```

### 2. Start MinIO with Docker Compose

The MinIO service is configured in `docker-compose.yaml`. Start it with:

```bash
docker-compose up -d
```

This will start:

- PostgreSQL on port 5432
- MinIO on port 9000 (API)
- MinIO Console on port 9001 (Web UI)

### 3. Access MinIO Console

Visit `http://localhost:9001` and login with:

- Username: `minioadmin` (or your configured `MINIO_ROOT_USER`)
- Password: `minioadmin` (or your configured `MINIO_ROOT_PASSWORD`)

## API Usage

### Create Menu Item with Image

**Endpoint:** `POST /v1/menu-items`

**Content-Type:** `multipart/form-data`

**Authentication:** Cookie-based (session token cookie must be present)

**Form Fields:**

- `name` (required): Menu item name
- `sku` (required): Stock keeping unit
- `price_baht` (required): Price in Thai Baht
- `category_id` (optional): Category UUID
- `active` (optional): "true" or "false"
- `image` (optional): Image file (JPG, JPEG, PNG, GIF, WEBP, max 5MB)

**Example using cURL:**

```bash
curl -X POST http://localhost:8080/v1/menu-items \
  -b "token=YOUR_SESSION_TOKEN" \
  -F "name=Pad Thai" \
  -F "sku=PADTHAI001" \
  -F "price_baht=120" \
  -F "category_id=550e8400-e29b-41d4-a716-446655440000" \
  -F "active=true" \
  -F "image=@/path/to/image.jpg"
```

**Note:** `-b "token=YOUR_SESSION_TOKEN"` sends the cookie. Get the token by logging in first.

**Example using Postman:**

1. Set method to POST
2. Set URL to `http://localhost:8080/v1/menu-items`
3. **Authentication:** Login first to get session cookie (cookie is automatically sent)
4. Go to "Body" tab
5. Select "form-data"
6. Add the following key-value pairs:
   - `name`: Pad Thai (Text)
   - `sku`: PADTHAI001 (Text)
   - `price_baht`: 120 (Text)
   - `category_id`: 550e8400-e29b-41d4-a716-446655440000 (Text)
   - `active`: true (Text)
   - `image`: (select "File" type and upload an image)

### Update Menu Item with Image

**Endpoint:** `PUT /v1/menu-items/:id`

**Content-Type:** `multipart/form-data`

**Authentication:** Cookie-based (session token cookie must be present)

**Form Fields:** (same as Create)

**Notes:**

- If you upload a new image, the old image will be automatically deleted from MinIO
- If you don't include an `image` field, the existing image will be preserved

**Example using cURL:**

```bash
curl -X PUT http://localhost:8080/v1/menu-items/550e8400-e29b-41d4-a716-446655440000 \
  -b "token=YOUR_SESSION_TOKEN" \
  -F "name=Pad Thai Special" \
  -F "sku=PADTHAI001" \
  -F "price_baht=150" \
  -F "category_id=550e8400-e29b-41d4-a716-446655440000" \
  -F "active=true" \
  -F "image=@/path/to/new-image.jpg"
```

### Delete Menu Item

**Endpoint:** `DELETE /v1/menu-items/:id`

**Notes:**

- Deleting a menu item will automatically delete its associated image from MinIO

**Example using cURL:**

```bash
curl -X DELETE http://localhost:8080/v1/menu-items/550e8400-e29b-41d4-a716-446655440000 \
  -b "token=YOUR_SESSION_TOKEN"
```

## Image Constraints

- **Allowed formats:** JPG, JPEG, PNG, GIF, WEBP
- **Maximum size:** 5 MB
- **Storage:** Images are stored in the MinIO bucket (default: `menu-images`)
- **Access:** Images are publicly accessible via their URL

## Response Format

When you create or update a menu item, the response includes the `image_url`:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Pad Thai",
  "sku": "PADTHAI001",
  "price_baht": 120,
  "active": true,
  "image_url": "http://localhost:9000/menu-images/550e8400-1234-5678-9012-abcdef123456-1697486400.jpg",
  "category": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "name": "Thai Dishes",
    "display_order": 1
  }
}
```

## Architecture

### Components

1. **MinIO Client** (`database/minio.go`)

   - Initializes connection to MinIO
   - Creates bucket if it doesn't exist
   - Sets public read policy for images

2. **Storage Utilities** (`utils/storage.go`)

   - `UploadImageToMinio()`: Handles image upload with validation
   - `DeleteImageFromMinio()`: Handles image deletion

3. **Handler** (`feature/menuItem/delivery/http.go`)
   - Modified to accept multipart/form-data
   - Handles image upload in Create/Update operations
   - Cleans up images on errors and deletions

### Data Flow

1. Client sends multipart form data with image
2. Handler validates form fields and extracts image
3. Image is uploaded to MinIO via utility function
4. MinIO returns public URL
5. URL is saved in database along with menu item data
6. On update: new image is uploaded, old image is deleted
7. On delete: menu item and associated image are both deleted

## Production Considerations

### Security

1. **Change default credentials:**

   ```bash
   MINIO_ROOT_USER=your-secure-username
   MINIO_ROOT_PASSWORD=your-secure-password-min-8-chars
   ```

2. **Enable SSL:**

   ```bash
   MINIO_USE_SSL=true
   ```

3. **Use environment-specific endpoints:**
   ```bash
   # Production
   MINIO_ENDPOINT=minio.yourdomain.com:443
   ```

### Performance

- Consider using a CDN in front of MinIO for better image delivery
- Implement image optimization/resizing before upload if needed
- Set appropriate cache headers for static assets

### Backup

- Configure MinIO bucket replication for disaster recovery
- Set up regular backups of MinIO data volume

## Troubleshooting

### MinIO Connection Failed

**Error:** `Failed to initialize MinIO client`

**Solution:**

1. Check if MinIO is running: `docker ps`
2. Verify environment variables are set correctly
3. Check MinIO logs: `docker logs minio`

### Image Upload Failed

**Error:** `Failed to upload file to MinIO`

**Solution:**

1. Check file size (max 5MB)
2. Verify file format is allowed
3. Check MinIO bucket permissions
4. Verify bucket exists (should be auto-created)

### Images Not Accessible

**Error:** Image URL returns 403 Forbidden

**Solution:**

1. Check bucket policy is set to public read
2. Verify the bucket policy in MinIO Console
3. Re-run the application to recreate the bucket with correct policy

## Testing

### Manual Testing with Postman

See the example requests above and configure Postman collection accordingly.

### Automated Testing

You can test the endpoints using the existing test framework. Remember to:

1. Mock the MinIO client in unit tests
2. Use a separate test bucket for integration tests
3. Clean up test images after each test

## References

- [MinIO Documentation](https://min.io/docs/minio/linux/index.html)
- [MinIO Go SDK](https://github.com/minio/minio-go)
- [S3 API Reference](https://docs.aws.amazon.com/AmazonS3/latest/API/Welcome.html)
