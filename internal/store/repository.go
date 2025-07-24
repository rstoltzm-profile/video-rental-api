package store

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StoreReader interface {
	CountTitlesByStore(ctx context.Context, storeID int) ([]StoreInventorySummary, error)
}

type Repository interface {
	StoreReader
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

func (r *repository) CountTitlesByStore(ctx context.Context, storeID int) ([]StoreInventorySummary, error) {
	query := `
		SELECT
			store.store_id, film.title,
			COUNT(film.title) AS title_count
		FROM
			inventory
			INNER JOIN store ON inventory.store_id = store.store_id
			INNER JOIN film ON inventory.film_id = film.film_id
		WHERE
			store.store_id = $1
		GROUP BY store.store_id, film.title;
	`
	rows, err := r.pool.Query(ctx, query, storeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var storeInventorySummary []StoreInventorySummary
	for rows.Next() {
		var c StoreInventorySummary
		if err := rows.Scan(&c.StoreID, &c.Title, &c.TitleCount); err != nil {
			return nil, err
		}
		storeInventorySummary = append(storeInventorySummary, c)
	}
	return storeInventorySummary, nil
}
