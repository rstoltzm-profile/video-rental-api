## Added List Rentals By Customer and int testing
```
5. List Rentals by Customer
API:
GET /customers/{id}/rentals
GET /customers/{id}/rentals?status=late
Useful for showing customer rental history or late returns.
```

## POST /payments
```
6. Payment Integration (Basic)
Record payments made by customers.
API:
POST /payments
Payload: {rental_id, amount, customer_id, staff_id}
✅ Make payment succeeded: 32106
```

## payment table needs partitions
```
CREATE TABLE public.payment_p2025_07 PARTITION OF public.payment
FOR VALUES FROM ('2025-07-01 00:00:00+00') TO ('2025-08-01 00:00:00+00');
``` 