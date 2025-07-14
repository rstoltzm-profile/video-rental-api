# Customer Routes
* http://localhost:8080/v1/

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | /customers | Get all customers |
| GET | /customers/{id} | Get a customer by ID |
| POST | /customers | Create a new customer|
| DELETE | /customers/{id} | Delete customer by ID |

### Create Customer POST /customers
* Request
```json
{
  "store_id": 1,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com",
  "address": {
    "address": "123 Main St",
    "address2": "",
    "district": "Downtown",
    "city_name": "Chicago",
    "postal_code": "90210",
    "phone": "555-1234"
  }
}
```
* Response
```json
{
  "id": 601,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john@example.com"
}
```