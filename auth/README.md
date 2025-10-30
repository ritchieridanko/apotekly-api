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
├── cmd/
│  ├── app/
│  └── migrate/
├── configs/
├── internal/
│  ├── app/
│  │  ├── caches/
│  │  ├── publishers/
│  │  ├── repositories/
│  │  └── usecases/
│  ├── entities/
│  ├── infrastructure/
│  │  ├── broker/
│  │  ├── cache/
│  │  ├── database/
│  │  ├── logger/
│  │  └── tracer/
│  ├── interfaces/
│  │  ├── di/
│  │  └── http/
│  │     ├── dto/
│  │     ├── handlers/
│  │     ├── middlewares/
│  │     ├── router/
│  │     └── validator/
│  ├── servers/
│  ├── services/
│  │  ├── broker/
│  │  ├── cache/
│  │  ├── database/
│  │  ├── logger/
│  │  └── oauth/
│  │     ├── google/
│  │     └── microsoft/
│  └── shared/
│     ├── ce/
│     ├── constants/
│     └── utils/
├── migrations/
├── pkg/
│  └── events/
```

## 🚀 Running the Service

1. **Configure Environment**

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

2. **Run the Service**

   ```bash
   # This service depends on redis, kafka, and jaeger containers running, and its kafka topics registered.
   ```

   Build the docker images of the service and its dependencies

   ```bash
   make docker-build
   ```

   Run the docker containers

   ```bash
   make docker-up
   ```

   Apply database migrations

   ```bash
   make docker-migrate-up
   ```

3. **Stop the Service**

   Rollback database migrations

   ```bash
   make docker-migrate-down
   ```

   Stop the docker containers

   ```bash
   make docker-down
   ```

## 📖 API Endpoints

You can access the API documentation here: https://documenter.getpostman.com/view/33667018/2sB3HqHJDN
