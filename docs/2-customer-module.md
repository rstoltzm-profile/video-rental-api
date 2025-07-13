# 2 - Initial Customer Module

✅ Summary:
- Added customer handler, service, repository, and model.
- Implemented `/v1/customers` and `/v1/customers/{id}` endpoints.
- Used native Go 1.22 `http.ServeMux` routing with `PathValue`.

📁 Files Added:
- `internal/customer/handler.go`
- `internal/customer/service.go`
- `internal/customer/repository.go`
- `internal/customer/model.go`

## Get all customers
```bash
curl http://localhost:8080/v1/customers
[{"id":1,"first_name":"MARY","last_name":"SMITH","email":"MARY.SMITH@sakilacustomer.org"},{"id":2,"first_name":"PATRICIA","last_name":"JOHNSON","email":"PATRICIA.JOHNSON@sakilacustomer.org"},{"id":3,"first_name":"LINDA","last_name":"WILLIAMS","email":"LINDA.WILLIAMS@sakilacustomer.org"},{"id":4,"first_name":"BARBARA","last_name":"JONES","email":"BARBARA.JONES@sakilacustomer.org"},{"id":5,"first_name":"ELIZABETH","last_name":"BROWN","email":"ELIZABETH.BROWN@sakilacustomer.org"},{"id":6,"first_name":"JENNIFER","last_name":"DAVIS","email":"JENNIFER.DAVIS@sakilacustomer.org"},{"id":7,"first_name":"MARIA","last_name":"MILLER","email":"MARIA.MILLER@sakilacustomer.org"},{"id":8,"first_name":"SUSAN","last_name":"WILSON","email":"SUSAN.WILSON@sakilacustomer.org"},{"id":9,"first_name":"MARGARET","last_name":"MOORE","email":"MARGARET.MOORE@sakilacustomer.org"},
...]
```

## Get single customer
```bash
curl http://localhost:8080/v1/customers/25
{"id":25,"first_name":"DEBORAH","last_name":"WALKER","email":"DEBORAH.WALKER@sakilacustomer.org"}
curl http://localhost:8080/v1/customers/26
{"id":26,"first_name":"JESSICA","last_name":"HALL","email":"JESSICA.HALL@sakilacustomer.org"}
curl http://localhost:8080/v1/customers/27
{"id":27,"first_name":"SHIRLEY","last_name":"ALLEN","email":"SHIRLEY.ALLEN@sakilacustomer.org"}
```