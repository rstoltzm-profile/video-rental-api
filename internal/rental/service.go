package rental

import (
	"context"
)

type Service interface {
	GetRentals(ctx context.Context) ([]Rental, error)
	GetLateRentals(ctx context.Context) ([]Rental, error)
}

type service struct {
	reader RentalReader
	tx     TransactionManager
}

func NewService(reader RentalReader, tx TransactionManager) Service {
	return &service{
		reader: reader,
		tx:     tx,
	}
}

func (s *service) GetRentals(ctx context.Context) ([]Rental, error) {
	return s.reader.GetLateRentals(ctx)
}

func (s *service) GetLateRentals(ctx context.Context) ([]Rental, error) {
	return s.reader.GetLateRentals(ctx)
}
