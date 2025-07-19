package film

import (
	"context"
)

type Service interface {
	GetFilms(ctx context.Context) ([]Film, error)
	GetFilmByID(ctx context.Context, id int) (Film, error)
}

type service struct {
	reader FilmReader
	tx     TransactionManager
}

func NewService(reader FilmReader, tx TransactionManager) Service {
	return &service{
		reader: reader,
		tx:     tx,
	}
}

func (s *service) GetFilms(ctx context.Context) ([]Film, error) {
	return s.reader.GetFilms(ctx)
}

func (s *service) GetFilmByID(ctx context.Context, id int) (Film, error) {
	return s.reader.GetFilmByID(ctx, id)
}
