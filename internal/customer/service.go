package customer

import "context"

type Service interface {
	GetCustomers(ctx context.Context) ([]Customer, error)
	GetCustomerByID(ctx context.Context, id int) (Customer, error)
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
