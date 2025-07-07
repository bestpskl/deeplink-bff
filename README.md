# Openapi Service Template V2

This project implements a robust Go service following Clean Architecture principles and Domain-Driven Design (DDD) patterns. It provides a structured, maintainable, and scalable foundation for building microservices.

## Features

- Clean Architecture implementation with clear separation of concerns
- Domain-Driven Design (DDD) principles
- Middleware support for logging, recovery, and request/response handling
- Swagger API documentation
- Configurable logging system with censoring capabilities
- Session management
- Utility packages for common operations
- Makefile for easy project management

## Architecture Overview

The project follows a layered architecture based on Clean Architecture principles:

### Core Layer (`internal/core`)

The heart of the application containing the business logic and domain model:

- `domain/` - Core business entities and value objects
- `ports/` - Interface definitions for repositories and services
- `services/` - Implementation of business logic and use cases

### Adapters Layer (`internal/adapters`)

Handles communication between the core and external components:

- `handler/` - HTTP request handlers and DTOs
- `repositories/` - Database repository implementations
- `microservices/` - External service integrations

### Infrastructure Layer (`internal/infrastructure`)

Manages technical concerns:

- Database connections
- Redis integration
- External service clients

## Project Structure

```
.
├── cmd/                    # Application entrypoints
│   └── deeplink-api/      # Main application
├── config/                # Configuration management
├── constant/              # Global constants and error codes
├── docs/                  # Swagger API documentation
├── internal/              # Private application code
│   ├── adapters/         # External adapters
│   ├── infrastructure/   # Technical implementations
│   └── core/             # Business logic
├── middleware/           # HTTP middleware components
└── pkg/                  # Shared utilities
```

## Getting Started

1. Clone the repository

```bash
git clone https://gitdev.devops.krungthai.com/open-api/poc/openapi-service-template-v2.git
```

2. Install dependencies

```bash
go mod download
```

3. Configure the application

```bash
cp config/config.example.go config/config.go
# Edit config.go with your settings
```

4. Run the application

```bash
make run
```

## Development

### Prerequisites

- Go 1.23.4 or higher
- Make
- Docker (optional, for containerization)

### Available Make Commands

```bash
make build      # Build the application
make test       # Run tests
make lint       # Run linters
make swagger    # Generate Swagger documentation
make run        # Run the application locally
make docker     # Build Docker image
```

### Adding New Features

1. Define domain models in `internal/core/domain`
2. Create repository interfaces in `internal/core/ports`
3. Implement business logic in `internal/core/services`
4. Add HTTP handlers in `internal/adapters/handler`
5. Update repository implementations in `internal/adapters/repositories`

## API Documentation

API documentation is available via Swagger UI at `/swagger/index.html` when running the application in development mode.

To regenerate Swagger documentation:

```bash
make swagger
```

## Middleware

The application includes several middleware components:

- **Logger**: Request/response logging with customizable censoring
- **Recovery**: Panic recovery and error handling
- **Dump**: Request/response dumping for debugging
- **Session**: Session management and context propagation

## Logging

The project uses a custom logging package (`pkg/logx`) with features like:

- Log level management
- Sensitive data censoring
- Structured logging
- Blacklist functionality
- Build information inclusion
