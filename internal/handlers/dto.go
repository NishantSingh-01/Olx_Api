package handlers

import (
	"fmt"
	"strings"
	"time"
)

type ListingInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	City        string  `json:"city"`
	Status      string  `json:"status"`
}
type CreateListingResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type validationError struct {
	Field string
	Msg   string
}

func (v validationError) Error() string {
	return fmt.Sprintf("%s:%s", v.Field, v.Msg)
}
func (req ListingInput) Validate() error {
	if strings.TrimSpace(req.Title) == "" {
		return validationError{
			Field: "title",
			Msg:   "Title is required",
		}
	}
	if strings.TrimSpace(req.City) == "" {
		return validationError{
			Field: "city",
			Msg:   "City is required",
		}
	}
	if strings.TrimSpace(req.Description) == "" {
		return validationError{
			Field: "description",
			Msg:   "Description is required",
		}
	}
	if req.Price < 0 {
		return validationError{
			Field: "price",
			Msg:   "Price cannot be negative",
		}
	}
	return nil
}
