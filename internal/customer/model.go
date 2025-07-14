package customer

type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreateCustomerRequest struct {
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	Email     string       `json:"email"`
	StoreID   int          `json:"store_id"`
	Address   AddressInput `json:"address"`
}

type AddressInput struct {
	Address    string `json:"address"`
	Address2   string `json:"address2"`
	District   string `json:"district"`
	CityID     int    `json:"city_id"`
	PostalCode string `json:"postal_code"`
	Phone      string `json:"phone"`
}
