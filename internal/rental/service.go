package rental

import (
	"context"
	"fmt"
)

type Service interface {
	GetRentals(ctx context.Context) ([]Rental, error)
	GetLateRentals(ctx context.Context) ([]Rental, error)
	GetRentalsByCustomerID(ctx context.Context, customerID int) ([]Rental, error)
	GetLateRentalsByCustomerID(ctx context.Context, customerID int) ([]Rental, error)
	CreateRental(ctx context.Context, req CreateRentalRequest) (int, error)
	ReturnRentalByID(ctx context.Context, id int) error
}

type service struct {
	reader RentalReader
	writer RentalWriter
	tx     TransactionManager
}

func NewService(reader RentalReader, writer RentalWriter, tx TransactionManager) Service {
	return &service{
		reader: reader,
		writer: writer,
		tx:     tx,
	}
}

func (s *service) GetRentals(ctx context.Context) ([]Rental, error) {
	return s.reader.GetRentals(ctx)
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

func (s *service) CreateRental(ctx context.Context, req CreateRentalRequest) (int, error) {
	// Check if the inventory is already rented and not yet returned
	activeRental, err := s.reader.GetActiveRentalByInventoryID(ctx, req.InventoryID)
	if err != nil {
		return 0, fmt.Errorf("failed to check inventory availability, %w", err)
	}
	if activeRental != nil {
		return 0, fmt.Errorf("inventory is already rented out, %d", req.InventoryID)
	}

	return s.writer.InsertRental(ctx, req)
}

func (s *service) ReturnRentalByID(ctx context.Context, id int) error {
	return s.writer.UpdateRentalByID(ctx, id)
}
