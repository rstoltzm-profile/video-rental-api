# 11-rental-post

## Goals

```
3. Create Rental Transaction (Checkout)
Allow a customer to rent a movie.
API:
POST /rentals
Payload: {customer_id, inventory_id, staff_id}
You’ll insert into rental and maybe payment tables here.
```

## rental table didn't auto increment
```
SELECT setval('rental_rental_id_seq', (SELECT MAX(rental_id) FROM rental));
```