package customer

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockCustomerReader struct {
	mock.Mock
}

func (m *mockCustomerReader) GetAll(ctx context.Context) ([]Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Customer), args.Error(1)
}

func (m *mockCustomerReader) GetByID(ctx context.Context, id int) (Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Customer), args.Error(1)
}

func (m *mockCustomerReader) GetCityIDByName(ctx context.Context, tx pgx.Tx, cityName string) (int, error) {
	args := m.Called(ctx, tx, cityName)
	return args.Int(0), args.Error(1)
}

type mockWriter struct{}
type mockTxManager struct{}

func (m *mockWriter) InsertAddress(ctx context.Context, tx pgx.Tx, address AddressInput, cityID int) (int, error) {
	return 0, nil
}
func (m *mockWriter) InsertCustomer(ctx context.Context, tx pgx.Tx, req CreateCustomerRequest, addressID int) (*Customer, error) {
	return nil, nil
}
func (m *mockWriter) DeleteCustomerByID(ctx context.Context, id int) error {
	return nil
}
func (m *mockTxManager) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

func TestService_GetCustomerByID(t *testing.T) {
	// Step 1: Create a mock that implements CustomerReader
	mockReader := new(mockCustomerReader)

	// Step 2: Inject the mock into the service
	svc := NewService(mockReader, &mockWriter{}, &mockTxManager{})

	// Step 3: Define the expected return value for the mock
	expected := Customer{ID: 1, FirstName: "Test"}

	// Step 4: Tell the mock how to respond when GetByID is called
	mockReader.On("GetByID", mock.Anything, 1).Return(expected, nil)

	got, err := svc.GetCustomerByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	mockReader.AssertExpectations(t)
}
