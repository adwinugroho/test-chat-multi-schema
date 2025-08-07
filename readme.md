### Chat Multi Schema

This project is a multi-tenant chat API service, focusing on the `tenant-api` folder. It is designed to support multiple tenants with isolated data schemas, user management, and message handling.

#### Project Structure (tenant-api)

- **main.go**: Entry point of the application.
- **config/**: Configuration files for environment variables, database, and RabbitMQ.
- **controller/**: HTTP route handlers for tenants, users, and messages.
  - **middleware/**: Middleware for validation and request processing.
- **domain/**: Core domain models and interfaces for tenants, users, and messages.
- **model/**: Data transfer objects (DTOs) for requests and responses.
- **pkg/**: Shared utilities and helpers.
  - **helper/**: General helper functions.
  - **logger/**: Logging utilities.
  - **server/**: Server startup logic.
- **repository/**: Data access layer for tenants, users, and messages.
- **service/**: Business logic for tenants, users, messages, and background workers.

#### Demo Credentials

- Email: admin@example.com
- Password: password!
