# 11-rental-post

## Goals

```
3. Create Rental Transaction (Checkout)
Allow a customer to rent a movie.
API:
POST /rentals
Payload: {customer_id, inventory_id, staff_id}
You’ll insert into rental and maybe payment tables here.

4. Return Movie
Mark a rental as returned.
API:
POST /rentals/{id}/return
(Updates rental.return_date)
```

## rental table didn't auto increment
```
SELECT setval('rental_rental_id_seq', (SELECT MAX(rental_id) FROM rental));
```

## Added Rental and Return Rental Feature
```
	mux.HandleFunc("POST /rentals", handler.CreateRental)
	mux.HandleFunc("POST /rentals/{id}/return", handler.ReturnRental)
```

## Added validation to look for rentals that are already checked out
```
bash test/integration-rental.sh
## Rent a movie
✅ POST /v1/rentals ... OK (created rental ID = 16080)
✅ POST /v1/rentals/16080/return ... OK
✅ POST /v1/rentals (first attempt) ... OK (created rental ID = 16081)
✅ POST /v1/rentals (second attempt, should fail) ... OK (duplicate rental rejected with status 500)
```