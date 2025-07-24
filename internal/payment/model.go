package payment

type Payment struct {
	CustomerID int     `json:"customer_id" validate:"required,gt=0"`
	StaffID    int     `json:"staff_id" validate:"required,gt=0"`
	RentalID   int     `json:"rental_id" validate:"required,gt=0"`
	Amount     float32 `json:"amount" validate:"required"`
}
