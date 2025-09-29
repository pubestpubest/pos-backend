# Go Clean Architecture Template

A production-ready template for building scalable and maintainable Go applications using Clean Architecture principles.

## 🚀 Features

- Clean Architecture implementation
- RESTful API with Gin framework
- PostgreSQL database integration
- Environment-based configuration
- Structured logging with Logrus
- CORS middleware
- Health check endpoint
- Versioned API routes
- Feature-based organization
- Domain-driven design
- Standardized error handling

## 📁 Project Structure

```
.
├── configs/              # Configuration files
├── configs.example/      # Example configuration files
├── constant/            # Global constants
├── database/            # Database connection and migrations
├── domain/              # Core business logic and entities
├── feature/             # Feature modules
├── middlewares/         # HTTP middlewares
├── models/              # Data models
├── request/             # Request DTOs
├── response/            # Response DTOs
├── routes/              # API route definitions
├── utils/               # Utility functions
│   └── error.go         # Error handling utilities
├── main.go              # Application entry point
└── go.mod               # Go module definition
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- PostgreSQL
- Make (optional, for using Makefile commands)

## 🏁 Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/pubestpubest/go-clean-arch-template.git
   cd go-clean-arch-template
   ```

2. Copy example configuration:
   ```bash
   cp -r configs.example/* configs/
   ```

3. Update the configuration files in `configs/` with your environment-specific settings.

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`.

## 🔧 Configuration

The application uses environment variables for configuration. Create a `.env` file in the `configs/` directory with the following variables:

```env
RUN_ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=your_database
```

## 📚 API Documentation

### Health Check
- `GET /healthz` - Health check endpoint

### API Version 1
- Base URL: `/v1`

## 🔍 Logging

The application uses Logrus for structured logging with the following features:

### Log Levels
- `Info` - General operational entries about what's happening inside the application
- `Warn` - Warning messages that don't necessarily affect the application's operation
- `Error` - Error events that might still allow the application to continue running
- `Fatal` - Critical errors that force the application to exit

### Log Format
```go
log.SetFormatter(&log.TextFormatter{
    ForceColors:   true,
    FullTimestamp: true,
})
```

### Example Usage
```go
// Info level logging
log.Info("Application started successfully")

// Warning level logging
log.Warn("Resource usage is high")

// Error level logging
log.Error("Failed to connect to database")

// Fatal level logging (will exit the application)
log.Fatal("Critical error occurred")
```

### Environment-based Logging
- Development mode: Colorized output with full timestamps
- Production mode: Plain text output with essential information only

## 🚨 Error Handling

The template implements a standardized error handling approach using the `utils/error.go` utility.

### Error Wrapping
Errors are wrapped with context using the `errors.Wrap` function from the `github.com/pkg/errors` package:

```go
err = errors.Wrap(err, "[UserRepository.GetUser]: Error getting user")
```

### Standard Error Formatting
The `StandardError` function in `utils/error.go` is used to extract clean error messages:

```go
func StandardError(err error) string {
    errorMessages := strings.Split(err.Error(), "]: ")
    return errorMessages[len(errorMessages)-1]
}
```

### Usage in Handlers
In HTTP handlers, errors are processed and returned in a standardized format:

```go
func (h *UserHandler) GetUser(c *gin.Context) {
    user, err := h.userService.GetUser(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": utils.StandardError(err),
        })
        return
    }
    c.JSON(http.StatusOK, user)
}
```

This approach provides:
- Consistent error formatting across the application
- Clean error messages for API responses
- Detailed error context in logs
- Easy error tracking through the call stack

## 📚 Libraries and Dependencies

### Core Libraries
- [Gin](https://github.com/gin-gonic/gin) - High-performance HTTP web framework
- [Logrus](https://github.com/sirupsen/logrus) - Structured logging
- [godotenv](https://github.com/joho/godotenv) - Environment variable management
- [GORM](https://gorm.io/) - ORM for database operations
- [PostgreSQL](https://www.postgresql.org/) - Database system

### Development Tools
- [Go](https://golang.org/) - Programming language
- [Make](https://www.gnu.org/software/make/) - Build automation tool

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin) - For providing a robust and efficient web framework
- [Logrus](https://github.com/sirupsen/logrus) - For structured logging capabilities
- [godotenv](https://github.com/joho/godotenv) - For environment variable management
- [GORM](https://gorm.io/) - For database operations and migrations
- [PostgreSQL](https://www.postgresql.org/) - For reliable database management
