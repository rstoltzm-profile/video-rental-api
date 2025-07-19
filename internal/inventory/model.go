package inventory

import "time"

type Inventory struct {
	InventoryID int       `json:"inventory_id"`
	LastUpdate  time.Time `json:"last_update"`
	FilmID      int       `json:"film_id"`
	Title       string    `json:"title"`
	StoreID     int       `json:"store_id"`
	AddressId   int       `json:"address_id"`
	Phone       string    `json:"phone"`
}

type InventoryAvailability struct {
	InventoryID int    `json:"inventory_id"`
	StoreID     int    `json:"store_id"`
	FilmID      int    `json:"film_id"`
	Title       string `json:"title"`
}
