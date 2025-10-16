# Environment Variables Configuration

Create a `.env` file in the `configs/` directory with the following variables:

## Application Environment

```bash
RUN_ENV=development
DEPLOY_ENV=local
SEED_DB=false
```

## Database Configuration

```bash
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=postgres
DATABASE_NAME=pos_db
```

## MinIO Object Storage Configuration

```bash
# MinIO server endpoint (without http/https)
MINIO_ENDPOINT=localhost:9000

# MinIO credentials (default for development)
MINIO_ROOT_USER=minioadmin
MINIO_ROOT_PASSWORD=minioadmin

# SSL Configuration
MINIO_USE_SSL=false

# Default bucket name for menu item images
MINIO_BUCKET_NAME=menu-images

# Ports for docker-compose
MINIO_PORT=9000
MINIO_CONSOLE_PORT=9001
```

## JWT Configuration

```bash
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h
```

## Notes

- For production, make sure to change `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD` to secure values
- Set `MINIO_USE_SSL=true` in production
- The MinIO console will be accessible at `http://localhost:9001` (or the port you specify)
- Images will be stored in the bucket specified by `MINIO_BUCKET_NAME` (default: `menu-images`)
