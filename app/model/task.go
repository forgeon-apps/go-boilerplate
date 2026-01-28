package model

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	ProjectID uuid.UUID  `db:"project_id" json:"project_id" validate:"required"`
	Title     string     `db:"title" json:"title" validate:"required"`
	Status    string     `db:"status" json:"status"` // todo|doing|done (we’ll default to todo)
	DueAt     *time.Time `db:"due_at" json:"due_at"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
}
