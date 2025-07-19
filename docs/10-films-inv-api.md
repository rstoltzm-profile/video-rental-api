# 10-films-inv-api

✅ Summary:
- Added Inventory Available Check
- Added Available false, for no sql returned

```text
curl -s "$BASE_URL/v1/inventory/available?film_id=1&store_id=2"
{"inventory_id":7,"store_id":2,"film_id":1,"title":"ACADEMY DINOSAUR","available":true}

curl -s "$BASE_URL/v1/inventory/available?film_id=999999&store_id=2"
{"available":false,"film_id":999999,"store_id":2}

```

## I ran into a unique edge case with the database

## Bad SQL Query
* It shows the film_id has being out of stock
* The problem has to do with the rental table, which if a customer does not return it marks return as [null] for history
* In this case the query would pickup this null and show as unavailable.
```SQL
SELECT
	inv.inventory_id,
	inv.store_id,
	inv.film_id,
	f.title
FROM
	inventory inv
	INNER JOIN film f ON inv.film_id = f.film_id
WHERE
	inv.store_id = 1
	AND inv.film_id = 586
	AND inv.inventory_id = 2670
	AND NOT EXISTS (
		SELECT 1
		FROM rental r
		WHERE r.inventory_id = inv.inventory_id
		AND r.return_date IS NULL
	)
LIMIT 1;

```

## Final SQL Query - orders the return 
```SQL
WITH latest_rentals AS (
	SELECT DISTINCT ON (inv.inventory_id)
		inv.inventory_id,
		inv.store_id,
		inv.film_id,
		f.title,
		r.rental_date,
		r.return_date
	FROM
		inventory inv
	JOIN film f ON inv.film_id = f.film_id
	LEFT JOIN rental r ON inv.inventory_id = r.inventory_id
	WHERE
		inv.store_id = 1
		AND inv.film_id = 586
	ORDER BY
		inv.inventory_id,
		r.rental_date DESC
)
SELECT *
FROM latest_rentals
WHERE return_date IS NOT NULL
LIMIT 1;
```

## Refactored a bit, since I didn't want to return the rental info here, and want TRUE AS available
```
	WITH latest_rentals AS (
		SELECT DISTINCT ON (inv.inventory_id)
			inv.inventory_id,
			inv.store_id,
			inv.film_id,
			f.title,
			r.rental_date,
			r.return_date
		FROM
			inventory inv
		JOIN film f ON inv.film_id = f.film_id
		LEFT JOIN rental r ON inv.inventory_id = r.inventory_id
		WHERE
			inv.store_id = $1
			AND inv.film_id = $2
		ORDER BY
			inv.inventory_id,
			r.rental_date DESC
	)
	SELECT 
		inventory_id,
		store_id,
		film_id,
		title,
		TRUE AS available
	FROM latest_rentals
	WHERE return_date IS NOT NULL
	LIMIT 1;
```