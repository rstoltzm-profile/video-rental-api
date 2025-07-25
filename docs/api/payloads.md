# Payload Examples

## Customer
```json
{
  "store_id": 1,
  "first_name": "Good",
  "last_name": "Doe",
  "email": "john@example.com",
  "address": {
    "address": "123 Main St",
    "address2": "",
    "district": "Sind",
    "city_name": "Shikarpur",
    "postal_code": "90210",
    "phone": "+15551234"
  }
}
```

## rental
```json
{
  "inventory_id": 709,
  "customer_id": 397,
  "staff_id": 1
}
```

## payment
```json
{
  "customer_id": 1,
  "staff_id": 1,
  "rental_id": 1,
  "amount": 0.99
}
```

## login
```json
{
  "username":"staff1", "password":"password123"
}
```

