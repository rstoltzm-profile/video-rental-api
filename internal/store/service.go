package store

import (
	"context"
)

type Service interface {
	GetStoreInventorySummary(ctx context.Context, storeID int) ([]StoreInventorySummary, error)
}

type service struct {
	reader StoreReader
	tx     TransactionManager
}

func NewService(reader StoreReader, tx TransactionManager) Service {
	return &service{
		reader: reader,
		tx:     tx,
	}
}

func (s *service) GetStoreInventorySummary(ctx context.Context, storeID int) ([]StoreInventorySummary, error) {
	return s.reader.CountTitlesByStore(ctx, storeID)
}
