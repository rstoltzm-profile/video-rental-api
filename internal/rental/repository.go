package rental

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type RentalReader interface {
	GetRentals(ctx context.Context) ([]Rental, error)
	GetLateRentals(ctx context.Context) ([]Rental, error)
}

type Repository interface {
	RentalReader
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
	rows, err := r.conn.Query(ctx, query)
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
	rows, err := r.conn.Query(ctx, query)
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
