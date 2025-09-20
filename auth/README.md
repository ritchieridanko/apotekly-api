# Apotekly - Auth Service

The **Auth Service** is responsible for handling authentication and authorization in the Apotekly platform. This service provides features like:

- User Registration
- JWT-based Authentication
- Session Management
- OAuth Integration with Google and Microsoft
- Email Verification
- Password Resets

## 📂 Project Structure

```bash
auth/
├── cmd/             # Application entrypoints
│  ├── migrate/      # DB migration application
│  └── server/       # Main server application
├── config/          # Configuration management
├── database/        # DB migrations and seedings
│  └── migrations/   # SQL migration files
└── internal/        # Private application code
   ├── ce/           # Custom error handlers
   ├── constants/    # Application constants
   ├── di/           # Dependency injection container
   ├── dto/          # Data transfer objects
   ├── entities/     # Domain entities
   ├── handlers/     # HTTP request handlers
   ├── infras/       # Infrastructure initializations
   ├── middlewares/  # HTTP middlewares
   ├── repos/        # Data access layer
   ├── routers/      # Route definitions
   ├── server/       # Server setup and configuration
   ├── services/     # Business logic services
   │  ├── cache/     # Caching service layer
   │  ├── db/        # DB service layer
   │  ├── email/     # Email service with templates
   │  └── oauth/     # OAuth provider integrations
   ├── usecases/     # Usecase implementations
   ├── utils/        # Utility functions
   └── validators/   # Custom validation rules
```

## 🚀 Running the Service

1. **Configure Environment**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Start Dependencies**

   Make sure you have PostgreSQL and Redis installed on your device. If not, follow the installation guidelines below:

   - `PostgreSQL`: https://www.postgresql.org/download/
   - `Redis`: https://redis.io/docs/latest/operate/oss_and_stack/install/archive/install-redis/

   Alternatively, if you have Docker on your device, you can install these dependencies from your Docker. Here's how to install the dependencies with Docker:

   ```bash
   docker run -d --name apotekly-postgres \
      -e POSTGRES_USER=postgres \
      -e POSTGRES_PASSWORD=postgres \
      -p 5432:5432 postgres:15

   docker run -d --name apotekly-redis \
      -p 6379:6379 redis:latest
   ```

3. **Run the Service**

   You can run the service in development mode

   ```bash
   make dev-up
   ```

   Or, you can build and run the service

   ```bash
   make build-and-run
   ```

## 📖 API Endpoints

You can access the API documentation here: https://documenter.getpostman.com/view/33667018/2sB3HqHJDN
