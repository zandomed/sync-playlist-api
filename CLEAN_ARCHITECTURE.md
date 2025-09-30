# Clean Architecture Implementation

This document describes the Clean Architecture refactoring of the sync-playlist-api project, following Robert C. Martin's Clean Architecture principles and Clean Code practices.

## Architecture Overview

The project now follows a layered architecture with clear dependency inversion:

```
┌─────────────────────────────────────────────────────────────┐
│                    External Interfaces                     │
│  (Web, Database, External APIs, File System, etc.)        │
├─────────────────────────────────────────────────────────────┤
│                     Interface Adapters                     │
│     (Controllers, Gateways, Presenters, etc.)             │
├─────────────────────────────────────────────────────────────┤
│                      Use Cases                            │
│           (Application Business Rules)                     │
├─────────────────────────────────────────────────────────────┤
│                      Entities                             │
│          (Enterprise Business Rules)                       │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
internal/
├── domain/                    # Domain Layer (innermost)
│   ├── entities/             # Business entities
│   │   ├── user.go
│   │   └── account.go
│   ├── valueobjects/         # Value objects
│   │   ├── user_id.go
│   │   ├── email.go
│   │   ├── user_profile.go
│   │   ├── account_id.go
│   │   └── hashed_password.go
│   ├── repositories/         # Repository interfaces
│   │   ├── user_repository.go
│   │   ├── account_repository.go
│   │   └── token_repository.go
│   └── errors/               # Domain-specific errors
│       └── domain_error.go
│
├── usecases/                 # Use Cases Layer
│   └── auth/
│       ├── register_user.go
│       └── login_user.go
│
└── adapters/                 # Adapters Layer (outermost)
    ├── repositories/         # Repository implementations
    │   ├── postgres_user_repository.go
    │   ├── postgres_account_repository.go
    │   └── postgres_token_repository.go
    ├── auth/                 # Authentication adapters
    │   └── jwt_token_generator.go
    ├── http/                 # HTTP adapters
    │   ├── handlers/
    │   │   └── auth_handler.go
    │   ├── dtos/
    │   │   └── auth_dtos.go
    │   ├── mappers/
    │   │   └── auth_mapper.go
    │   └── routes/
    │       └── routes.go
    ├── container/            # Dependency injection
    │   └── container.go
    └── database/
        └── migrations/
            └── 002_add_refresh_tokens_table.sql
```

## Layer Descriptions

### 1. Domain Layer (Core)

The innermost layer containing the business logic and rules.

#### Entities
- **User**: Core user entity with business validation
- **Account**: Authentication accounts (userpass, OAuth providers)

#### Value Objects
- **UserID**: Type-safe user identifier
- **Email**: Email with validation
- **UserProfile**: User profile information
- **HashedPassword**: Secure password handling

#### Repository Interfaces
- Abstract interfaces defining data access contracts
- No implementation details, only business contracts

#### Domain Errors
- Custom error types with business context
- Better error handling than generic errors

### 2. Use Cases Layer (Application)

Contains application-specific business rules and orchestrates domain entities.

#### Features
- **RegisterUser**: User registration with validation
- **LoginUser**: User authentication with token generation
- Clean separation of concerns
- No dependencies on external frameworks

### 3. Adapters Layer (Infrastructure)

Implements interfaces from inner layers and handles external concerns.

#### Repository Implementations
- PostgreSQL implementations of repository interfaces
- Proper error handling and mapping

#### HTTP Adapters
- Controllers, DTOs, mappers
- Framework-specific logic isolated
- Clean request/response handling

#### Authentication Adapters
- JWT token generation and validation
- Configurable and replaceable

## Key Clean Architecture Principles Applied

### 1. Dependency Inversion
- Inner layers define interfaces
- Outer layers implement those interfaces
- Dependencies point inward only

### 2. Clean Code Practices
- **Meaningful Names**: Clear, descriptive variable and function names
- **Small Functions**: Single responsibility, easy to understand
- **SOLID Principles**:
  - Single Responsibility: Each class has one reason to change
  - Open/Closed: Open for extension, closed for modification
  - Liskov Substitution: Subtypes must be substitutable for base types
  - Interface Segregation: Many specific interfaces better than one general
  - Dependency Inversion: Depend on abstractions, not concretions

### 3. Separation of Concerns
- Business logic separated from infrastructure
- Web framework details isolated
- Database implementation details hidden

### 4. Testability
- Each layer can be tested independently
- Mock implementations for interfaces
- No external dependencies in business logic

## Usage

### Running with Clean Architecture

```go
// Use the new main file
go run cmd/server/main_clean.go
```

### Configuration

The system uses the same configuration but with updated JWT settings:

```env
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
JWT_REFRESH_EXPIRATION=8640000  # 100 days in seconds
```

### API Endpoints

The API maintains the same interface but with improved error handling:

```
POST /api/v1/auth/register
POST /api/v1/auth/login
```

## Benefits of This Architecture

1. **Maintainability**: Clear separation makes code easier to understand and modify
2. **Testability**: Each layer can be tested in isolation
3. **Flexibility**: Easy to swap implementations (e.g., different databases)
4. **Scalability**: Well-organized code scales better
5. **Domain Focus**: Business logic is protected and clearly defined
6. **Error Handling**: Proper domain-specific error handling
7. **Type Safety**: Strong typing with value objects prevents errors

## Migration Guide

### From Old to New Architecture

1. **Entities**: Old models → New domain entities with behavior
2. **Services**: Old services → Use cases + domain services
3. **Repositories**: Old repositories → Interface + implementation separation
4. **Handlers**: Old handlers → Clean handlers with proper DTOs
5. **Errors**: Generic errors → Domain-specific errors

### Gradual Migration

This refactoring can be done gradually:
1. Start with domain layer
2. Add use cases
3. Create adapters
4. Update handlers
5. Switch dependency injection

## Testing Strategy

### Unit Tests
- Domain entities and value objects
- Use cases with mocked dependencies
- Repository implementations

### Integration Tests
- HTTP endpoints
- Database operations
- Full authentication flow

### Example Test Structure
```go
// Domain entity test
func TestUser_ChangeEmail(t *testing.T) { /* ... */ }

// Use case test with mocks
func TestRegisterUser_Success(t *testing.T) { /* ... */ }

// Integration test
func TestAuthHandler_Register_Success(t *testing.T) { /* ... */ }
```

This Clean Architecture implementation provides a solid foundation for maintainable, testable, and scalable Go applications.