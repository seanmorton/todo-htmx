package domain

import (
	"time"
)

type Project struct {
	Id        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (p Project) IsDeleted() bool {
	return p.DeletedAt != nil
}

type ProjectFilters struct {
	ShowDeleted bool
}
