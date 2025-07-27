

# Video Rental API - Go Development Guide

## Changelog Index

- See changes marked as x-name.md in 
  - docs/1-mvp-build-process/
  - docs/2-mvp-post-process/
- Number increments as project moves on

## Architecture & Structure
- **Layered Domain Architecture**: Each domain follows `handler → service → repository` pattern
- **Domain Packages**: `/internal/{customer,rental,inventory,store,film}` - each contains handler, service, repository, model
- **Interface-Driven Design**: Services depend on Reader/Writer interfaces for testability
- **Single DB Connection**: Uses `*pgx.Conn` injected into all repositories (no connection pooling yet)

## Key Patterns

### Service Registration (internal/api/router.go)
```go
repo := customer.NewRepository(conn)
svc := customer.NewService(repo, repo, repo) // reader, writer, tx interfaces
handler := customer.NewHandler(svc)
```

### Transaction Pattern (customer/service.go)
Complex operations use explicit transactions with proper rollback:
```go
tx, err := s.tx.BeginTx(ctx)
defer tx.Rollback(ctx)
// ... multiple operations
tx.Commit(ctx)
```

### Interface Segregation
- Separate `Reader` and `Writer` interfaces per domain  
- `TransactionManager` interface for tx operations
- Services inject all needed interfaces, not concrete repositories

## Development Workflow

### Build & Test
```bash
make build          # Builds to bin/video-rental-api
make test           # Unit tests
make integration-test # Runs test/integration-*.sh scripts
make run            # Local development
```

### Database Setup
Requires Pagila sample database (DVD rental schema):
```bash
export DATABASE_URL="postgres://postgres:123456@localhost:5432/postgres"
# Alternative port if conflicts: localhost:6543
```

### Test Data
- `test/customer.json` - Customer creation payload
- `test/rental.json` - Rental payload (inventory_id: 709, customer_id: 397)
- Integration tests expect specific IDs (customer 373, store 1, film 1)

## Route Patterns

### RESTful Structure
- Health: `/health`  
- API: `/v1/{domain}` with proper HTTP methods
- Path params: `r.PathValue("id")` (Go 1.22+ routing)
- Query params: `r.URL.Query().Get("late")`, `customer_id`, etc.

### Error Handling
Currently basic `http.Error()` responses. Missing:
- Structured error types
- Proper status codes based on error type
- Request validation

## Testing Approach

### Unit Tests (internal/*/service_test.go)
- Mock interfaces using testify/mock
- Test business logic in services
- Example: `customer/service_test.go` shows mock pattern

### Integration Tests (test/*.sh)
- Bash scripts using curl + jq
- Test full HTTP endpoints with real database
- Expects server running on localhost:8080

## Current Limitations
- No input validation (struct tags exist but not enforced)
- Single DB connection (should use connection pool)
- Basic error responses (need structured errors)
- No middleware (logging, auth, CORS)
- No observability

## External Dependencies
- `github.com/jackc/pgx/v5` - PostgreSQL driver
- `github.com/stretchr/testify` - Testing utilities
- Pagila database schema (external dependency)

## Module Path
`github.com/rstoltzm-profile/video-rental-api`