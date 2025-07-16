package store

type StoreInventorySummary struct {
	StoreID    int    `json:"store_id"`
	Title      string `json:"title"`
	TitleCount int    `json:"title_count"`
}
