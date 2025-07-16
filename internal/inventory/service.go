package inventory

import (
	"context"
)

type Service interface {
	GetInventory(ctx context.Context) ([]Inventory, error)
	GetInventoryByStore(ctx context.Context, storeID int) ([]Inventory, error)
}

type service struct {
	reader InventoryReader
	tx     TransactionManager
}

func NewService(reader InventoryReader, tx TransactionManager) Service {
	return &service{
		reader: reader,
		tx:     tx,
	}
}

func (s *service) GetInventory(ctx context.Context) ([]Inventory, error) {
	return s.reader.GetInventory(ctx)
}

func (s *service) GetInventoryByStore(ctx context.Context, storeID int) ([]Inventory, error) {
	return s.reader.GetInventoryByStore(ctx, storeID)
}
