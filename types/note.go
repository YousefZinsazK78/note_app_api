package types

type Note struct {
	ID          int    `json:"id,omitempty" form:"id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}
