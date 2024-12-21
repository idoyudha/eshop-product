package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        string
	Name      string
	ParentID  *string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (p *Category) GenerateCategoryID() error {
	categoryID, err := uuid.NewV7()
	if err != nil {
		return err
	}

	p.ID = categoryID.String()
	return nil
}
