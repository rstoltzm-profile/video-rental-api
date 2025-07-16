package inventory

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type InventoryReader interface {
	GetInventory(ctx context.Context) ([]Inventory, error)
	GetInventoryByStore(ctx context.Context, storeID int) ([]Inventory, error)
}

type Repository interface {
	InventoryReader
	TransactionManager
}

type TransactionManager interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type repository struct {
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return &repository{conn: conn}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.conn.Begin(ctx)
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
	rows, err := r.conn.Query(ctx, query)
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
	rows, err := r.conn.Query(ctx, query, storeID)
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
