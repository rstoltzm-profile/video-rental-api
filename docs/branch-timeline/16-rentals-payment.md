## Added List Rentals By Customer and int testing
```
5. List Rentals by Customer
API:
GET /customers/{id}/rentals
GET /customers/{id}/rentals?status=late
Useful for showing customer rental history or late returns.
```

## ToDO
```
6. Payment Integration (Basic)
Record payments made by customers.
API:
POST /payments
Payload: {rental_id, amount, customer_id, staff_id}
Uses the payment table.
```