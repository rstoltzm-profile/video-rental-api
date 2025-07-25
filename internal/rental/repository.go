package rental

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RentalReader interface {
	GetRentals(ctx context.Context) ([]Rental, error)
	GetLateRentals(ctx context.Context) ([]Rental, error)
	GetActiveRentalByInventoryID(ctx context.Context, inventoryID int) (*Rental, error)
}

type RentalWriter interface {
	InsertRental(ctx context.Context, req CreateRentalRequest) (int, error)
	UpdateRentalByID(ctx context.Context, id int) error
}

type Repository interface {
	RentalReader
	TransactionManager
	RentalWriter
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

func (r *repository) GetRentals(ctx context.Context) ([]Rental, error) {
	query := `
	SELECT
		customer.first_name,
		customer.last_name, 
		address.phone,
		rental.rental_date,
		film.title
	FROM
		rental
		INNER JOIN customer ON rental.customer_id = customer.customer_id
		INNER JOIN address ON customer.address_id = address.address_id
		INNER JOIN inventory ON rental.inventory_id = inventory.inventory_id
		INNER JOIN film ON inventory.film_id = film.film_id
	WHERE
		rental.return_date IS NULL
	ORDER BY
		rental.rental_date
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []Rental
	for rows.Next() {
		var c Rental
		if err := rows.Scan(&c.FirstName, &c.LastName, &c.Phone, &c.RentalDate, &c.Title); err != nil {
			return nil, err
		}
		rentals = append(rentals, c)
	}
	return rentals, nil
}

func (r *repository) GetLateRentals(ctx context.Context) ([]Rental, error) {
	query := `
	SELECT
		customer.first_name,
		customer.last_name, 
		address.phone,
		rental.rental_date,
		film.title
	FROM
		rental
		INNER JOIN customer ON rental.customer_id = customer.customer_id
		INNER JOIN address ON customer.address_id = address.address_id
		INNER JOIN inventory ON rental.inventory_id = inventory.inventory_id
		INNER JOIN film ON inventory.film_id = film.film_id
	WHERE
		rental.return_date IS NULL
		AND rental_date < CURRENT_DATE
	ORDER BY
		rental.rental_date
	`
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rentals []Rental
	for rows.Next() {
		var c Rental
		if err := rows.Scan(&c.FirstName, &c.LastName, &c.Phone, &c.RentalDate, &c.Title); err != nil {
			return nil, err
		}
		rentals = append(rentals, c)
	}
	return rentals, nil
}

func (r *repository) InsertRental(ctx context.Context, req CreateRentalRequest) (int, error) {
	var rental_id int
	query := `
		INSERT INTO rental (rental_date, inventory_id, customer_id, staff_id, last_update)
		VALUES (
		TO_TIMESTAMP(TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD HH24:MI:SS'), 'YYYY-MM-DD HH24:MI:SS'),
		$1, $2, $3,
		TO_TIMESTAMP(TO_CHAR(CURRENT_TIMESTAMP, 'YYYY-MM-DD HH24:MI:SS'), 'YYYY-MM-DD HH24:MI:SS')
		)
		RETURNING rental_id
	`
	err := r.pool.QueryRow(ctx, query, req.InventoryID, req.CustomerID, req.StaffID).Scan(&rental_id)

	if err != nil {
		return -1, err
	}

	return rental_id, nil
}

func (r *repository) UpdateRentalByID(ctx context.Context, id int) error {
	query := `
	UPDATE rental
	SET return_date = CURRENT_TIMESTAMP
	WHERE rental_id = $1
	`
	cmdTag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("update rental failed: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no rental found with ID %d", id)
	}
	return nil
}

func (r *repository) GetActiveRentalByInventoryID(ctx context.Context, inventoryID int) (*Rental, error) {
	var rental Rental
	query := `
	SELECT
		customer.first_name,
		customer.last_name, 
		address.phone,
		rental.rental_date,
		film.title
	FROM
		rental
		INNER JOIN customer ON rental.customer_id = customer.customer_id
		INNER JOIN address ON customer.address_id = address.address_id
		INNER JOIN inventory ON rental.inventory_id = inventory.inventory_id
		INNER JOIN film ON inventory.film_id = film.film_id
	WHERE
		inventory.inventory_id = $1
		and rental.return_date IS NULL
	ORDER BY
		rental.rental_date
	`
	err := r.pool.QueryRow(ctx, query, inventoryID).Scan(&rental.FirstName, &rental.LastName, &rental.Phone, &rental.RentalDate, &rental.Title)
	if err == pgx.ErrNoRows {
		return nil, nil // No Active rental
	}
	if err != nil {
		return nil, err
	}

	return &rental, err
}
