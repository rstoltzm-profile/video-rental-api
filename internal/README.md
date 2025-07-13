# modular design
```
internal/customer/
├── handler.go         # HTTP layer
├── service.go         # Business logic
├── repository.go      # DB access
├── model.go           # Domain types (optional, or shared in internal/model)
```

## api version
```
managed in internal/api/router.go

later we can add v2
example:
internal/customer/
    v2/handler.go
    v2/service.go
```