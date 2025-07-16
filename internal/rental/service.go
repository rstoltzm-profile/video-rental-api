package rental

import (
	"context"
)

type Service interface {
	GetRentals(ctx context.Context) ([]Rental, error)
	GetLateRentals(ctx context.Context) ([]Rental, error)
	GetRentalsByCustomerID(ctx context.Context, customerID int) ([]Rental, error)
	GetLateRentalsByCustomerID(ctx context.Context, customerID int) ([]Rental, error)
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

func (s *service) GetRentalsByCustomerID(ctx context.Context, customerID int) ([]Rental, error) {
	return s.reader.GetRentalsByCustomerID(ctx, customerID)
}

func (s *service) GetLateRentals(ctx context.Context) ([]Rental, error) {
	return s.reader.GetLateRentals(ctx)
}

func (s *service) GetLateRentalsByCustomerID(ctx context.Context, customerID int) ([]Rental, error) {
	return s.reader.GetLateRentalsByCustomerID(ctx, customerID)
}
