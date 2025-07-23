package customer

type Customer struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type CreateCustomerRequest struct {
	FirstName string       `json:"first_name" validate:"required,min=1,max=50"`
	LastName  string       `json:"last_name" validate:"required,min=1,max=50"`
	Email     string       `json:"email" validate:"required,email"`
	StoreID   int          `json:"store_id" validate:"required,gt=0"`
	Address   AddressInput `json:"address" validate:"required"`
}

type AddressInput struct {
	Address    string `json:"address" validate:"required,min=1,max=100"`
	Address2   string `json:"address2" validate:"max=100"`
	District   string `json:"district" validate:"required,min=1,max=50"`
	CityName   string `json:"city_name" validate:"required,min=1,max=50"`
	PostalCode string `json:"postal_code" validate:"required,min=4,max=6"`
	Phone      string `json:"phone" validate:"required,e164"`
}
