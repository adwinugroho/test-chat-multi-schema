### Chat Multi Schema

This project is a multi-tenant chat API service, focusing on the `tenant-api` folder. It is designed to support multiple tenants with isolated data schemas, user management, and message handling.

---

#### Example Configuration

Copy and edit `tenant-api/config/example.config.yaml` as needed:

```yaml
app_version: "0.0.1"
app_name: "tenant-api"
app_port: 9000
app_url: "http://localhost"
environment: "development"
jwt_secret: ""
rabbitmq:
  url: amqp://user:pass@localhost:5672/
database:
  url: postgres://user:pass@localhost:5432/app
workers: 3 # Default worker
```

---

#### Database Schema Example

```sql
CREATE TABLE tenants (
    tenant_id UUID PRIMARY KEY,
    tenant_name TEXT,
    user_id UUID NOT NULL
);

CREATE TABLE users (
    user_id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password TEXT,
    role TEXT,
    tenant_id UUID
);

CREATE TABLE messages (
    message_id UUID,
    tenant_id UUID NOT NULL,
    payload JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (tenant_id, message_id)
) PARTITION BY LIST (tenant_id);
```

---

#### Usage Flow

1. **Insert a user** (example):
   - Insert a user into the `users` table for login.
2. **Login**:
   - Use the login endpoint to get a JWT token.
   - Copy the response token and add it to the `Authorization` header for requests to the `/tenants` endpoint.
3. **Create a tenant**:
   - Use the `/tenants` endpoint to create a tenant. Save the `tenant_id` from the response.
4. **Access messages**:
   - For requests to the `/messages` endpoint, add the `X-Tenant-ID` header with the value of the `tenant_id` you created.

---

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

#### TO DO

- [x] retry logic dead letter
- [x] generate swagger documentation
- [x] add more integration test
- [x] add service background to kill all or specific consumer
- [x] add service tenant manager (middleware service)
- [x] add SSE endpoint with redis for realtime
- [x] Implement Prometheus monitoring
