package film

import (
	"context"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const baseFilmQuery = `
	SELECT title, description, release_year, language.name, rating
	FROM film
	INNER JOIN language on film.language_id = language.language_id
`

const baseFilmWithActorsQuery = `
	SELECT 
	  film.film_id,
	  film.title, 
	  film.description, 
	  film.release_year, 
	  language.name AS language, 
	  film.rating, 
	  category.name AS category,
	  actor.first_name || ' ' || actor.last_name AS actor_name
	FROM film
	INNER JOIN language ON film.language_id = language.language_id
	INNER JOIN film_category ON film.film_id = film_category.film_id
	INNER JOIN category ON film_category.category_id = category.category_id
	INNER JOIN film_actor ON film.film_id = film_actor.film_id
	INNER JOIN actor ON film_actor.actor_id = actor.actor_id
	WHERE film.film_id = $1
`

type FilmReader interface {
	GetFilms(ctx context.Context) ([]Film, error)
	GetFilmByID(ctx context.Context, id int) (Film, error)
	FindByTitle(ctx context.Context, title string) ([]Film, error)
	FindFilmWithActorsAndCategoriesByID(ctx context.Context, id int) (FilmWithActorsCategories, error)
}

type Repository interface {
	FilmReader
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

func (r *repository) GetFilms(ctx context.Context) ([]Film, error) {
	query := baseFilmQuery
	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var c Film
		if err := rows.Scan(&c.Title, &c.Description, &c.ReleaseYear, &c.Language, &c.Rating); err != nil {
			return nil, err
		}
		films = append(films, c)
	}
	return films, nil
}

func (r *repository) GetFilmByID(ctx context.Context, id int) (Film, error) {
	var c Film
	query := baseFilmQuery + ` WHERE film.film_id = $1`

	err := r.pool.QueryRow(ctx, query, id).Scan(&c.Title, &c.Description, &c.ReleaseYear, &c.Language, &c.Rating)

	return c, err
}

func (r *repository) FindByTitle(ctx context.Context, title string) ([]Film, error) {
	query := baseFilmQuery + ` WHERE film.title = $1`
	rows, err := r.pool.Query(ctx, query, title)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var films []Film
	for rows.Next() {
		var c Film
		if err := rows.Scan(&c.Title, &c.Description, &c.ReleaseYear, &c.Language, &c.Rating); err != nil {
			return nil, err
		}
		films = append(films, c)
	}
	return films, nil

}

func (r *repository) FindFilmWithActorsAndCategoriesByID(ctx context.Context, id int) (FilmWithActorsCategories, error) {
	rows, err := r.pool.Query(ctx, baseFilmWithActorsQuery, id)
	if err != nil {
		return FilmWithActorsCategories{}, err
	}
	defer rows.Close()

	var film FilmWithActorsCategories
	firstRow := true
	categoriesMap := make(map[string]struct{})
	actorsMap := make(map[string]struct{})

	for rows.Next() {
		var filmID int
		var actor string
		var category string

		err := rows.Scan(
			&filmID,
			&film.Title,
			&film.Description,
			&film.ReleaseYear,
			&film.Language,
			&film.Rating,
			&category,
			&actor,
		)
		if err != nil {
			return FilmWithActorsCategories{}, err
		}

		if firstRow {
			firstRow = false
		}
		categoriesMap[category] = struct{}{}
		actorsMap[actor] = struct{}{}
	}

	if firstRow {
		// no rows found
		return FilmWithActorsCategories{}, pgx.ErrNoRows
	}

	// unique categories list
	categories := make([]string, 0, len(categoriesMap))
	sort.Strings(categories)
	for c := range categoriesMap {
		categories = append(categories, c)
	}
	film.Categories = categories

	// unique actors list
	actors := make([]string, 0, len(actorsMap))
	sort.Strings(actors)
	for c := range actorsMap {
		actors = append(actors, c)
	}
	film.Actors = actors

	// trim whitespace
	film.Language = strings.TrimSpace(film.Language)

	return film, nil
}
