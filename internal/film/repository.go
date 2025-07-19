package film

import (
	"context"

	"github.com/jackc/pgx/v5"
)

const baseFilmQuery = `
	SELECT title, description, release_year, language.name, rating, category.name as category
	FROM film
	INNER JOIN language on film.language_id = language.language_id
	INNER JOIN film_category on film.film_id = film_category.film_id
	INNER JOIN category on film_category.category_id = category.category_id
`

type FilmReader interface {
	GetFilms(ctx context.Context) ([]Film, error)
	GetFilmByID(ctx context.Context, id int) (Film, error)
	FindByTitle(ctx context.Context, title string) ([]Film, error)
}

type Repository interface {
	FilmReader
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

func (r *repository) GetFilms(ctx context.Context) ([]Film, error) {
	query := baseFilmQuery
	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var c Film
		if err := rows.Scan(&c.Title, &c.Description, &c.ReleaseYear, &c.Language, &c.Rating, &c.Category); err != nil {
			return nil, err
		}
		films = append(films, c)
	}
	return films, nil
}

func (r *repository) GetFilmByID(ctx context.Context, id int) (Film, error) {
	var c Film
	query := baseFilmQuery + ` WHERE film.film_id = $1`

	err := r.conn.QueryRow(ctx, query, id).Scan(&c.Title, &c.Description, &c.ReleaseYear, &c.Language, &c.Rating, &c.Category)

	return c, err
}

func (r *repository) FindByTitle(ctx context.Context, title string) ([]Film, error) {
	query := baseFilmQuery + ` WHERE film.title = $1`
	rows, err := r.conn.Query(ctx, query, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var c Film
		if err := rows.Scan(&c.Title, &c.Description, &c.ReleaseYear, &c.Language, &c.Rating, &c.Category); err != nil {
			return nil, err
		}
		films = append(films, c)
	}
	return films, nil

}
