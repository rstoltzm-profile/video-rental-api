package rental

import "time"

type Rental struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Phone      string    `json:"phone"`
	RentalDate time.Time `json:"rental_date"`
	Title      string    `json:"title"`
}

type CreateRentalRequest struct {
	InventoryID int `json:"inventory_id" validate:"min=0"`
	CustomerID  int `json:"customer_id" validate:"min=0"`
	StaffID     int `json:"staff_id" validate:"min=0"`
}
