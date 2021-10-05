package domain

import (
	"time"
)

type Recipe struct {
	ID          string       `json:"id"`
	RecipeName  string       `json:"recipe_name"`
	Description string       `json:"description"`
	Ingredients []Ingredient `json:"ingredients"`
	OwnerID     string       `json:"owner_id"`
	CreatedAt   time.Time    `json:"created_at"`
	DeletedAt   time.Time    `json:"deleted_at"`
}

type Ingredient struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Quantity    float32 `json:"quantity"`
	Measurement string  `json:"measurement"`
}
