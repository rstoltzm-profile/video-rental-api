package customer

import (
	"context"
	"fmt"
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

func (m *mockCustomerReader) FindCustomerRentalsByID(ctx context.Context, id int) ([]CustomerRentals, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]CustomerRentals), args.Error(1)
}

func (m *mockCustomerReader) FindLateCustomerRentalsByID(ctx context.Context, id int) ([]CustomerRentals, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]CustomerRentals), args.Error(1)
}

func (m *mockCustomerReader) GetCityIDByName(ctx context.Context, cityName string) (int, error) {
	args := m.Called(ctx, cityName)
	return args.Int(0), args.Error(1)
}

type mockWriter struct {
	mock.Mock
}
type mockTxManager struct {
	mock.Mock
}

func (m *mockWriter) InsertAddress(ctx context.Context, address AddressInput, cityID int) (int, error) {
	args := m.Called(ctx, address, cityID)
	return args.Int(0), args.Error(1)
}

func (m *mockWriter) InsertCustomer(ctx context.Context, req CreateCustomerRequest, addressID int) (*Customer, error) {
	args := m.Called(ctx, req, addressID)

	var customer *Customer
	if c := args.Get(0); c != nil {
		customer = c.(*Customer)
	}
	return customer, args.Error(1)
}

func (m *mockWriter) DeleteCustomerByID(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockTxManager) BeginTx(ctx context.Context) (pgx.Tx, error) {
	args := m.Called(ctx)
	return args.Get(0).(pgx.Tx), args.Error(1)
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

func TestService_GetAll(t *testing.T) {
	mockReader := new(mockCustomerReader)
	svc := NewService(mockReader, &mockWriter{}, &mockTxManager{})

	expected := []Customer{{ID: 1, FirstName: "Test"}}
	mockReader.On("GetAll", mock.Anything).Return(expected, nil)

	got, err := svc.GetCustomers(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	mockReader.AssertExpectations(t)
}

func TestService_GetCustomerRentalsByID(t *testing.T) {
	mockReader := new(mockCustomerReader)
	svc := NewService(mockReader, &mockWriter{}, &mockTxManager{})

	expected := []CustomerRentals{{FirstName: "John", LastName: "House"}}
	mockReader.On("FindCustomerRentalsByID", mock.Anything, 1).Return(expected, nil)

	got, err := svc.GetCustomerRentalsByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	mockReader.AssertExpectations(t)
}

func TestService_GetLateCustomerRentalsByIDD(t *testing.T) {
	mockReader := new(mockCustomerReader)
	svc := NewService(mockReader, &mockWriter{}, &mockTxManager{})

	expected := []CustomerRentals{{FirstName: "John", LastName: "House"}}
	mockReader.On("FindLateCustomerRentalsByID", mock.Anything, 1).Return(expected, nil)

	got, err := svc.GetLateCustomerRentalsByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, expected, got)

	mockReader.AssertExpectations(t)
}

func TestService_DeleteCustomerByID(t *testing.T) {
	mockWriter := new(mockWriter)
	svc := NewService(&mockCustomerReader{}, mockWriter, &mockTxManager{})

	mockWriter.On("DeleteCustomerByID", mock.Anything, 1).Return(nil)

	err := svc.DeleteCustomerByID(context.Background(), 1)

	assert.NoError(t, err)

	mockWriter.AssertExpectations(t)
}

func TestService_DeleteCustomerByID_NotFound(t *testing.T) {
	mockWriter := new(mockWriter)
	svc := NewService(&mockCustomerReader{}, mockWriter, &mockTxManager{})

	mockWriter.On("DeleteCustomerByID", mock.Anything, 999).Return(fmt.Errorf("no customer found with ID 999"))

	err := svc.DeleteCustomerByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no customer found")
	mockWriter.AssertExpectations(t)
}

func TestService_CreateCustomer(t *testing.T) {
	mockReader := new(mockCustomerReader)
	mockWriter := new(mockWriter)

	svc := NewService(mockReader, mockWriter, nil) // no txManager needed now

	newCustomer := CreateCustomerRequest{
		StoreID:   1,
		FirstName: "John",
		LastName:  "House",
		Email:     "JohnHouse@rental.com",
		Address: AddressInput{
			// fill fields if your service uses any for cityName
			CityName: "TestCity",
		},
	}

	expectedCustomer := &Customer{
		ID:        1,
		FirstName: "John",
		LastName:  "House",
		Email:     "JohnHouse@rental.com",
	}

	// Mock GetCityIDByName to return a cityID (e.g., 42)
	mockReader.On("GetCityIDByName", mock.Anything, "TestCity").Return(42, nil)

	// Mock InsertAddress to return an addressID (e.g., 10)
	mockWriter.On("InsertAddress", mock.Anything, newCustomer.Address, 42).Return(10, nil)

	// Mock InsertCustomer to return expected customer
	mockWriter.On("InsertCustomer", mock.Anything, newCustomer, 10).Return(expectedCustomer, nil)

	resp, err := svc.CreateCustomer(context.Background(), newCustomer)

	assert.NoError(t, err)
	assert.Equal(t, expectedCustomer, resp)

	mockReader.AssertExpectations(t)
	mockWriter.AssertExpectations(t)
}
