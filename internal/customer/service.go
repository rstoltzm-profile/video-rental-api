package customer

import (
	"context"
	"fmt"
)

type Service interface {
	GetCustomers(ctx context.Context) ([]Customer, error)
	GetCustomerByID(ctx context.Context, id int) (Customer, error)
	CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*Customer, error)
	DeleteCustomerByID(ctx context.Context, id int) error
}

type service struct {
	reader CustomerReader
	writer CustomerWriter
	tx     TransactionManager
}

func NewService(reader CustomerReader, writer CustomerWriter, tx TransactionManager) Service {
	return &service{
		reader: reader,
		writer: writer,
		tx:     tx,
	}
}

func (s *service) GetCustomers(ctx context.Context) ([]Customer, error) {
	return s.reader.GetAll(ctx)
}

func (s *service) GetCustomerByID(ctx context.Context, id int) (Customer, error) {
	customer, err := s.reader.GetByID(ctx, id)
	if err != nil {
		return Customer{}, fmt.Errorf("Customer not found")
	}
	return customer, nil
}

func (s *service) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*Customer, error) {
	tx, err := s.tx.BeginTx(ctx) // helper to get pgx.Tx
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get CityID
	cityID, err := s.reader.GetCityIDByName(ctx, tx, req.Address.CityName)
	if err != nil {
		return nil, fmt.Errorf("invalid city name: %w", err)
	}

	// Insert Address
	addressID, err := s.writer.InsertAddress(ctx, tx, req.Address, cityID)
	if err != nil {
		return nil, err
	}

	// Insert Customer
	customer, err := s.writer.InsertCustomer(ctx, tx, req, addressID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *service) DeleteCustomerByID(ctx context.Context, id int) error {
	return s.writer.DeleteCustomerByID(ctx, id)
}
