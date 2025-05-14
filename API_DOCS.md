
# API Endpoints Documentation

## Base URL
All API endpoints are prefixed with `/api/v1/`

## Authentication and Authorization
- Most endpoints are protected by Keycloak/Okta authentication
- Roles-based access control is implemented using guards
- Available roles: `admin`, `appdev`, `dataops`, `devops`

## User Microservice

### Authentication Endpoints (`/authLocal`)
- `POST /authLocal` - Create a local admin user
- `PUT /authLocal/:id` - Edit a local admin user

### Users Endpoints (`/users`)
- `POST /users` - Create a new user
- `PUT /users/:id` - Edit an existing user
- `GET /users` - Get all users (Query params: createdBy)
- `GET /users/:id` - Get user by ID

### Roles Endpoints (`/roles`)
- `GET /roles/auth` - Check if user's token is valid (Public)
- `POST /roles` - Create a new role (Public)
- `PUT /roles/:id` - Edit an existing role (Public)
- `GET /roles` - Get all roles (Public)

### Admin Settings Endpoints (`/adminSettings`)
- `GET /adminSettings` - Get all admin settings (Public)

## Enbuild Microservice

### Operations Endpoints (`/operations`)
- `POST /operations` - Create a new operation
- `PUT /operations/:id` - Edit an existing operation
- `GET /operations` - Get all operations with filtering, pagination, and sorting
- `GET /operations/:id` - Get operation by ID

### Repository Endpoints (`/repository`)
- `GET /repository` - Get all repos
- `GET /repository/:id` - Get repo by ID

### Manifests Endpoints (`/manifests`)
- `GET /manifests` - Get all manifests
- `GET /manifests/:id` - Get manifest by ID

## ML Microservice

### ML Dataset Endpoints (`/mlDataset`)
- `GET /mlDataset` - Get all ML datasets

## Common Features

### Query Parameters
- Many GET endpoints support the following query parameters:
  - `limit` (default: 10) - Number of items per page
  - `page` (default: 1) - Page number
  - `sort` (default: '-createdOn') - Sort field and direction
  - Additional filters can be passed as query parameters

### Response Format
```json
{
  "data": {
    // Response data
  }
}
```

### Error Response Format
```json
{
  "statusCode": number,
  "message": "Error",
  "error": "Error message"
}
```

### Authentication Headers
- Bearer token authentication is required for protected endpoints
- Format: `Authorization: Bearer <token>`

### Version Control
- All endpoints are versioned (currently v1)
- Version is specified in the URL path

### Correlation ID
- All requests are tracked with a correlation ID for debugging and logging purposes
- Correlation ID is managed by the `CorrelationService`

This documentation covers the main API endpoints available in the backend microservices. Each endpoint is protected by appropriate authentication and authorization mechanisms, and follows RESTful conventions for request/response handling.
