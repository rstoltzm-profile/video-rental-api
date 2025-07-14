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
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetCustomers(ctx context.Context) ([]Customer, error) {
	return s.repo.GetAll(ctx)
}

func (s *service) GetCustomerByID(ctx context.Context, id int) (Customer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) CreateCustomer(ctx context.Context, req CreateCustomerRequest) (*Customer, error) {
	tx, err := s.repo.BeginTx(ctx) // helper to get pgx.Tx
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get CityID
	cityID, err := s.repo.GetCityIDByName(ctx, tx, req.Address.CityName)
	if err != nil {
		return nil, fmt.Errorf("invalid city name: %w", err)
	}

	// Insert Address
	addressID, err := s.repo.InsertAddress(ctx, tx, req.Address, cityID)
	if err != nil {
		return nil, err
	}

	// Insert Customer
	customer, err := s.repo.InsertCustomer(ctx, tx, req, addressID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return customer, nil
}

func (s *service) DeleteCustomerByID(ctx context.Context, id int) error {
	return s.repo.DeleteCustomerByID(ctx, id)
}
