package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/idoyudha/eshop-product/internal/utils"
)

type Product struct {
	ID          string
	SKU         string
	Name        string
	ImageURL    string
	Description string
	Price       float64
	Quantity    int
	CategoryID  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

func (p *Product) GenerateProductID() error {
	productId, err := uuid.NewV7()
	if err != nil {
		return err
	}

	p.ID = productId.String()
	return nil
}

func (p *Product) GenerateSKU() {
	p.SKU = utils.GenerateSKU()
}

func (p *Product) SetImageURL(imageURL string) {
	p.ImageURL = imageURL
}
