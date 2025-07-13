package customer

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetAll(ctx context.Context) ([]Customer, error)
	GetByID(ctx context.Context, id int) (Customer, error)
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
