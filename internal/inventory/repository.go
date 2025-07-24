package inventory

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryReader interface {
	GetInventory(ctx context.Context) ([]Inventory, error)
	GetInventoryByStore(ctx context.Context, storeID int) ([]Inventory, error)
	FindInventoryAvailable(ctx context.Context, storeID int, filmID int) (InventoryAvailability, error)
}

type Repository interface {
	InventoryReader
	TransactionManager
}

type TransactionManager interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) Repository {
	return &repository{pool: pool}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *repository) GetInventory(ctx context.Context) ([]Inventory, error) {
	query := `
	SELECT 	inventory.inventory_id,
			inventory.last_update,
			inventory.film_id,
			film.title,
			inventory.store_id,
			store.address_id,
			address.phone
	FROM
		inventory
		INNER JOIN store ON inventory.store_id = store.store_id
		INNER JOIN film ON inventory.film_id = film.film_id
		INNER JOIN address ON store.address_id = address.address_id
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []Inventory
	for rows.Next() {
		var c Inventory
		if err := rows.Scan(&c.InventoryID, &c.LastUpdate, &c.FilmID, &c.Title, &c.StoreID, &c.AddressId, &c.Phone); err != nil {
			return nil, err
		}
		rentals = append(rentals, c)
	}
	return rentals, nil
}

func (r *repository) GetInventoryByStore(ctx context.Context, storeID int) ([]Inventory, error) {
	query := `
	SELECT 	inventory.inventory_id,
			inventory.last_update,
			inventory.film_id,
			film.title,
			inventory.store_id,
			store.address_id,
			address.phone
	FROM
		inventory
		INNER JOIN store ON inventory.store_id = store.store_id
		INNER JOIN film ON inventory.film_id = film.film_id
		INNER JOIN address ON store.address_id = address.address_id
	WHERE
		store.store_id = $1
	`
	rows, err := r.pool.Query(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []Inventory
	for rows.Next() {
		var c Inventory
		if err := rows.Scan(&c.InventoryID, &c.LastUpdate, &c.FilmID, &c.Title, &c.StoreID, &c.AddressId, &c.Phone); err != nil {
			return nil, err
		}
		rentals = append(rentals, c)
	}
	return rentals, nil
}

func (r *repository) FindInventoryAvailable(ctx context.Context, storeID int, filmID int) (InventoryAvailability, error) {
	var i InventoryAvailability
	query := `
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
	`

	err := r.pool.QueryRow(ctx, query, storeID, filmID).Scan(&i.InventoryID, &i.StoreID, &i.FilmID, &i.Title, &i.Available)

	return i, err
}
