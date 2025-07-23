package errors

type CustomerError struct {
	Code    string
	Message string
	Cause   error
}

var (
	ErrCustomerNotFound    = &CustomerError{Code: "CUSTOMER_NOT_FOUND", Message: "customer not found"}
	ErrInvalidCustomerData = &CustomerError{Code: "INVALID_DATA", Message: "invalid customer data"}
)
