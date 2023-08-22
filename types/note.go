package types

import "fmt"

type Note struct {
	ID          int    `json:"id,omitempty" form:"id"`
	Title       string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
}

func (n Note) ValidateNote() error {
	if len(n.Title) == 0 {
		return fmt.Errorf("title must not be empty! %d", len(n.Title))
	}
	if len(n.Description) == 0 {
		return fmt.Errorf("description must not be empty! %d", len(n.Description))
	}
	return nil
}
