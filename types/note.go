package types

type Note struct {
	ID          int    `json:"-"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
