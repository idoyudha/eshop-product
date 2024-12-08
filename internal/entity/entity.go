package entity

import "time"

type Category struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	ParentID  *int       `json:"parent_id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
