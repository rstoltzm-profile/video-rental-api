package customer

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CustomerReader interface {
	GetAll(ctx context.Context) ([]Customer, error)
	GetByID(ctx context.Context, id int) (Customer, error)
	GetCityIDByName(ctx context.Context, tx pgx.Tx, cityName string) (int, error)
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
	conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) Repository {
	return &repository{conn: conn}
}

func (r *repository) GetAll(ctx context.Context) ([]Customer, error) {
	rows, err := r.conn.Query(ctx, `SELECT customer_id, first_name, last_name, email FROM customer`)
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
	err := r.conn.QueryRow(ctx,
		`SELECT customer_id, first_name, last_name, email FROM customer WHERE customer_id = $1`,
		id,
	).Scan(&c.ID, &c.FirstName, &c.LastName, &c.Email)

	return c, err
}

func (r *repository) GetCityIDByName(ctx context.Context, tx pgx.Tx, cityName string) (int, error) {
	var cityID int
	err := r.conn.QueryRow(ctx,
		`SELECT city_id FROM city where city = $1`,
		cityName,
	).Scan(&cityID)

	if err != nil {
		return 0, err
	}

	return cityID, err
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.conn.Begin(ctx)
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
	cmdTag, err := r.conn.Exec(ctx,
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
