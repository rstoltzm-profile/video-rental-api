package film

type Film struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ReleaseYear int    `json:"release_year"`
	Language    string `json:"language"`
	Rating      string `json:"rating"`
	Category    string `json:"category"`
}
