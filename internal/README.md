# modular design
```
internal/customer/
├── handler.go         # HTTP layer
├── service.go         # Business logic
├── repository.go      # DB access
├── model.go           # Domain types (optional, or shared in internal/model)
```