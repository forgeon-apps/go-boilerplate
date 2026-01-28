package model

import (
	"time"

	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Name         string    `db:"name" json:"name"`
	Username     *string   `db:"username" json:"username,omitempty"`
	PasswordHash *string   `db:"password_hash" json:"-"` // never expose
	IsActive     bool      `db:"is_active" json:"is_active"`
	IsAdmin      bool      `db:"is_admin" json:"is_admin"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	IsDeleted    bool      `db:"is_deleted"`
	UserName     string    `db:"username"`
	Password     string    `db:"password"`
	FirstName    string    `db:"first_name"`
	LastName     string    `db:"last_name"`
}

func NewUser() *User {
	return &User{}
}

type CreateUser struct {
	IsAdmin   bool   `json:"is_admin"`
	IsActive  bool   `json:"is_active"`
	UserName  string `json:"username" validate:"required,lte=50,gte=5"`
	Email     string `json:"email" validate:"required,email,lte=150"`
	Password  string `json:"password" validate:"required,lte=100,gte=10"`
	FirstName string `json:"first_name" validate:"required,lte=100"`
	LastName  string `json:"last_name" validate:"required,lte=100"`
}

type UpdateUser struct {
	IsAdmin   bool   `json:"is_admin"`
	IsActive  bool   `json:"is_active"`
	FirstName string `json:"first_name" validate:"required,lte=100"`
	LastName  string `json:"last_name" validate:"required,lte=100"`
}
