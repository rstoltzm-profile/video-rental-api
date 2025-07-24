package customer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CustomerReader interface {
	GetAll(ctx context.Context) ([]Customer, error)
	GetByID(ctx context.Context, id int) (Customer, error)
	GetCityIDByName(ctx context.Context, tx pgx.Tx, cityName string) (int, error)
	FindCustomerRentalsByID(ctx context.Context, id int) ([]CustomerRentals, error)
	FindLateCustomerRentalsByID(ctx context.Context, id int) ([]CustomerRentals, error)
}

type CustomerWriter interface {
	InsertAddress(ctx context.Context, tx pgx.Tx, address AddressInput, cityID int) (int, error)
	InsertCustomer(ctx context.Context, tx pgx.Tx, req CreateCustomerRequest, addressID int) (*Customer, error)
	DeleteCustomerByID(ctx context.Context, id int) error
}

type Repository interface {
	CustomerReader
	CustomerWriter
	BeginTx(ctx context.Context) (pgx.Tx, error)
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

func (r *repository) GetAll(ctx context.Context) ([]Customer, error) {
	rows, err := r.pool.Query(ctx, `SELECT customer_id, first_name, last_name, email FROM customer`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		var c Customer
		if err := rows.Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email); err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}
	return customers, nil
}

func (r *repository) GetByID(ctx context.Context, id int) (Customer, error) {
	var c Customer
	err := r.pool.QueryRow(ctx,
		`SELECT customer_id, first_name, last_name, email FROM customer WHERE customer_id = $1`,
		id,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email)

	return c, err
}

func (r *repository) GetCityIDByName(ctx context.Context, tx pgx.Tx, cityName string) (int, error) {
	var cityID int
	err := r.pool.QueryRow(ctx,
		`SELECT city_id FROM city where city = $1`,
		cityName,
	).Scan(&cityID)

	if err != nil {
		return 0, err
	}

	return cityID, err
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pool.Begin(ctx)
}

func (r *repository) InsertAddress(ctx context.Context, tx pgx.Tx, address AddressInput, cityID int) (int, error) {
	var id int
	err := tx.QueryRow(ctx, `
	INSERT INTO address (address, address2, district, city_id, postal_code, phone)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING address_id
	`,
		address.Address,
		address.Address2,
		address.District,
		cityID,
		address.PostalCode,
		address.Phone,
	).Scan(&id)
	return id, err
}

func (r *repository) InsertCustomer(ctx context.Context, tx pgx.Tx, req CreateCustomerRequest, addressID int) (*Customer, error) {
	var id int
	err := tx.QueryRow(ctx, `
		INSERT INTO customer (store_id, first_name, last_name, email, address_id, activebool, create_date, last_update, active)
		VALUES ($1, $2, $3, $4, $5, TRUE, CURRENT_DATE, CURRENT_TIMESTAMP, 1)
		RETURNING customer_id
	`,
		req.StoreID,
		req.FirstName,
		req.LastName,
		req.Email,
		addressID,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return &Customer{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}, nil
}

func (r *repository) DeleteCustomerByID(ctx context.Context, id int) error {
	cmdTag, err := r.pool.Exec(ctx,
		`DELETE FROM customer WHERE customer_id = $1`,
		id,
	)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no customer found with ID %d", id)
	}
	return nil
}

func (r *repository) FindCustomerRentalsByID(ctx context.Context, id int) ([]CustomerRentals, error) {
	query := `
	SELECT
		customer.first_name,
		customer.last_name, 
		address.phone,
		rental.rental_date,
		film.title,
		rental.rental_date + (film.rental_duration || ' days')::interval AS rental_due_date,
		CURRENT_DATE > rental.rental_date + (film.rental_duration || ' days')::interval as overdue
	FROM
		rental
		INNER JOIN customer ON rental.customer_id = customer.customer_id
		INNER JOIN address ON customer.address_id = address.address_id
		INNER JOIN inventory ON rental.inventory_id = inventory.inventory_id
		INNER JOIN film ON inventory.film_id = film.film_id
	WHERE
		rental.return_date IS NULL
		and customer.customer_id = $1
	ORDER BY
		rental.rental_date
	`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerRentals []CustomerRentals
	for rows.Next() {
		var c CustomerRentals
		if err := rows.Scan(&c.FirstName, &c.LastName, &c.Phone, &c.RentalDate, &c.Title, &c.RentalDueDate, &c.Overdue); err != nil {
			return nil, err
		}
		customerRentals = append(customerRentals, c)
	}
	fmt.Println(customerRentals)
	return customerRentals, nil
}

func (r *repository) FindLateCustomerRentalsByID(ctx context.Context, id int) ([]CustomerRentals, error) {
	query := `
	SELECT
		customer.first_name,
		customer.last_name, 
		address.phone,
		rental.rental_date,
		film.title,
		rental.rental_date + (film.rental_duration || ' days')::interval AS rental_due_date,
		CURRENT_DATE > rental.rental_date + (film.rental_duration || ' days')::interval as overdue
	FROM
		rental
		INNER JOIN customer ON rental.customer_id = customer.customer_id
		INNER JOIN address ON customer.address_id = address.address_id
		INNER JOIN inventory ON rental.inventory_id = inventory.inventory_id
		INNER JOIN film ON inventory.film_id = film.film_id
	WHERE
		rental.return_date IS NULL
		and customer.customer_id = $1
		and CURRENT_DATE > rental.rental_date + (film.rental_duration || ' days')::interval
	ORDER BY
		rental.rental_date
	`
	rows, err := r.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerRentals []CustomerRentals
	for rows.Next() {
		var c CustomerRentals
		if err := rows.Scan(&c.FirstName, &c.LastName, &c.Phone, &c.RentalDate, &c.Title, &c.RentalDueDate, &c.Overdue); err != nil {
			return nil, err
		}
		customerRentals = append(customerRentals, c)
	}
	fmt.Println(customerRentals)
	return customerRentals, nil
}
