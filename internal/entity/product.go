package entity

import "time"

type Product struct {
	ID          string     `json:"id"`
	SKU         string     `json:"sku"`
	Name        string     `json:"name"`
	ImageURL    string     `json:"image_url"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Quantity    int        `json:"quantity"`
	CategoryID  int        `json:"category_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}
