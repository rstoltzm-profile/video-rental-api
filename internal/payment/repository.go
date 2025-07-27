package payment

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentWriter interface {
	InsertPayment(ctx context.Context, req Payment) (int, error)
}

type TransactionManager interface {
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type Repository interface {
	PaymentWriter
	TransactionManager
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

func (r *repository) InsertPayment(ctx context.Context, req Payment) (int, error) {
	var payment_id int
	fmt.Println(req)
	query := `
		INSERT INTO payment (customer_id, staff_id, rental_id, amount, payment_date)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING payment_id
	`
	err := r.pool.QueryRow(ctx, query, req.CustomerID, req.StaffID, req.RentalID, req.Amount).Scan(&payment_id)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23514" { // check SQLSTATE code
				fmt.Printf("InsertPayment partition error: %v\n", pgErr.Message)
				return -1, fmt.Errorf("payment insert failed: partition missing for current date")
			}
		}
		fmt.Printf("InsertPayment error : %v\n", err)
		return -1, err
	}

	return payment_id, nil
}
