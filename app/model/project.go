package model

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID          uuid.UUID `db:"id" json:"id"`
	OwnerUserID uuid.UUID `db:"owner_user_id" json:"owner_user_id" validate:"required"`
	Name        string    `db:"name" json:"name" validate:"required"`
	Description *string   `db:"description" json:"description"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}
