package customer

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (m *mockRepo) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
}

func (m *mockRepo) GetByID(ctx context.Context, id int) (Customer, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(Customer), args.Error(1)
}

func (m *mockRepo) GetAll(ctx context.Context) ([]Customer, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Customer), args.Error(1)
}

func (m *mockRepo) DeleteCustomerByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(1)
}

func (m *mockRepo) GetCityIDByName(ctx context.Context, tx pgx.Tx, cityName string) (int, error) {
	args := m.Called(ctx, tx, cityName)
	return args.Int(0), args.Error(1)
}

func (m *mockRepo) InsertAddress(ctx context.Context, tx pgx.Tx, address AddressInput, cityID int) (int, error) {
	args := m.Called(ctx, tx, address, cityID)
	return args.Int(0), args.Error(1)
}

func (m *mockRepo) InsertCustomer(ctx context.Context, tx pgx.Tx, req CreateCustomerRequest, addressID int) (*Customer, error) {
	args := m.Called(ctx, tx, req, addressID)

	// Use type assertion to *Customer, not Customer
	var customer *Customer
	if args.Get(0) != nil {
		customer = args.Get(0).(*Customer)
	}

	return customer, args.Error(1)
}

func TestService_GetCustomerByID_Success(t *testing.T) {
	mockRepo := new(mockRepo)
	svc := NewService(mockRepo)

	expected := Customer{
		ID:        42,
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane@example.com",
	}

	mockRepo.On("GetByID", mock.Anything, 42).Return(expected, nil)

	got, err := svc.GetCustomerByID(context.Background(), 42)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)
	mockRepo.AssertExpectations(t)
}

func TestService_GetCityIDByName_Success(t *testing.T) {
	mockRepo := new(mockRepo)

	expectedCityID := 1
	cityName := "Chicago"

	mockRepo.On("GetCityIDByName", mock.Anything, mock.Anything, cityName).Return(expectedCityID, nil)

	got, err := mockRepo.GetCityIDByName(context.Background(), nil, cityName)

	assert.NoError(t, err)
	assert.Equal(t, expectedCityID, got)
	mockRepo.AssertExpectations(t)
}

// InsertAddress(ctx context.Context, tx pgx.Tx, address AddressInput, cityID int) (int, error)
type TestAddressInput struct {
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	District   string `json:"district"`
	CityName   string `json:"city_name"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}
