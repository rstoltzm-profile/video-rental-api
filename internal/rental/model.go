package rental

import "time"

type Rental struct {
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Phone      string    `json:"phone"`
	RentalDate time.Time `json:"rental_date"`
	Title      string    `json:"title"`
}
