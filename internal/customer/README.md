# Customer domain (handlers, service, repository)


- Create: internal/customer/
  - handler.go         → HTTP handler (exposes GET /customers)
  - service.go         → Business logic (optional layer)
  - repository.go      → DB access
  - model.go           → Customer struct
- Responsibility:
  - handler: decode/encode requests
  - service: apply filtering, etc.
  - repository: query customer table
  - model: define domain struct
