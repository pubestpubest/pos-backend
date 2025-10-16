# MinIO Object Storage Implementation Summary

## Overview

Successfully implemented MinIO object storage for menu item image uploads. Images are now stored in MinIO and their URLs are saved in the database.

## Changes Made

### 1. Docker Configuration

**File:** `docker-compose.yaml`

- Added MinIO service with health check
- Configured ports 9000 (API) and 9001 (Console)
- Added persistent volume for MinIO data

### 2. Database Layer

**File:** `database/minio.go` (NEW)

- Created MinIO client initialization
- Auto-creates `menu-images` bucket on startup
- Sets public read policy for image access
- Configurable via environment variables

### 3. Utility Functions

**File:** `utils/storage.go` (NEW)

- `UploadImageToMinio()`: Handles image upload with validation
  - Validates file type (JPG, JPEG, PNG, GIF, WEBP)
  - Validates file size (max 5MB)
  - Generates unique filenames
  - Returns public URL
- `DeleteImageFromMinio()`: Handles image deletion from storage

### 4. Handler Updates

**File:** `feature/menuItem/delivery/http.go`

- Updated handler to accept MinIO client dependency
- Modified `CreateMenuItem()`:
  - Changed from JSON to multipart/form-data
  - Handles image upload
  - Cleans up on error
- Modified `UpdateMenuItem()`:
  - Changed from JSON to multipart/form-data
  - Uploads new image if provided
  - Deletes old image when replaced
  - Preserves existing image if no new upload
- Modified `DeleteMenuItem()`:
  - Deletes associated image from MinIO

### 5. Route Configuration

**File:** `routes/menuItemRoute.go`

- Updated to pass MinIO client to handler

### 6. Application Initialization

**File:** `main.go`

- Added MinIO connection on startup
- Fatal error if MinIO connection fails

### 7. Documentation

**Files Created:**

- `configs.example/ENV_VARIABLES.md`: Environment variable documentation
- `MINIO_IMAGE_UPLOAD_GUIDE.md`: Comprehensive usage guide

### 8. Dependencies

**File:** `go.mod`

- Added `github.com/minio/minio-go/v7` dependency

## Required Environment Variables

Add these to your `configs/.env` file:

```bash
MINIO_ENDPOINT=localhost:9000
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin
MINIO_USE_SSL=false
MINIO_BUCKET_NAME=menu-images
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
```

## API Changes

### Before (JSON)

```bash
POST /v1/menu-items
Content-Type: application/json

{
  "name": "Pad Thai",
  "sku": "PADTHAI001",
  "price_baht": 120,
  "image_url": "manual-url-here"
}
```

### After (Multipart Form Data)

```bash
POST /v1/menu-items
Content-Type: multipart/form-data

name: Pad Thai
sku: PADTHAI001
price_baht: 120
image: [file upload]
```

## Features

✅ Image upload with validation (type, size)
✅ Automatic unique filename generation
✅ Public URL generation
✅ Automatic cleanup on errors
✅ Old image deletion on update
✅ Associated image deletion on menu item deletion
✅ Preserves existing image if no new upload
✅ Auto-creates bucket with public read policy
✅ Configurable via environment variables

## Testing

1. Start services:

   ```bash
   docker-compose up -d
   ```

2. Access MinIO Console:

   - URL: http://localhost:9001
   - Username: minioadmin
   - Password: minioadmin

3. Test image upload:
   ```bash
   curl -X POST http://localhost:8080/v1/menu-items \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -F "name=Test Item" \
     -F "sku=TEST001" \
     -F "price_baht=100" \
     -F "image=@/path/to/image.jpg"
   ```

## Next Steps

1. Start Docker services: `docker-compose up -d`
2. Add environment variables to `configs/.env`
3. Run the application
4. Test image upload using the API

## Architecture Compliance

✅ Follows Clean Architecture principles
✅ Separation of concerns maintained
✅ Repository pattern preserved
✅ Error handling follows project standards
✅ Logging implemented consistently
✅ Infrastructure layer properly separated

## Production Considerations

- Change MinIO credentials for production
- Enable SSL (MINIO_USE_SSL=true)
- Consider CDN for image delivery
- Implement image optimization if needed
- Set up bucket replication for backups
- Monitor storage usage
