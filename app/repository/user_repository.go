package repository

import (
	"fmt"
	"time"

	"github.com/byeblogs/go-boilerplate/app/model"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/google/uuid"
)

type UserRepo struct {
	db *database.DB
}

func NewUserRepo(db *database.DB) UserRepository {
	return &UserRepo{db}
}

func (repo *UserRepo) Create(u *model.User) error {
	query := `
		INSERT INTO users (id, email, name, username, password_hash, is_active, is_admin, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7, now(), now())
	`
	_, err := repo.db.Exec(query, u.ID, u.Email, u.Name, u.Username, u.PasswordHash, u.IsActive, u.IsAdmin)
	return err
}

func (repo *UserRepo) All(limit int, offset uint) ([]*model.User, error) {
	var out []*model.User
	query := `SELECT * FROM users ORDER BY created_at DESC`
	var err error

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
		err = repo.db.Select(&out, query, limit, offset)
	} else {
		err = repo.db.Select(&out, query)
	}
	return out, err
}

func (repo *UserRepo) Get(id uuid.UUID) (*model.User, error) {
	u := model.User{}
	query := `SELECT * FROM users WHERE id = $1`
	if err := repo.db.Get(&u, query, id); err != nil {
		return nil, err
	}
	return &u, nil
}

func (repo *UserRepo) GetByUsername(username string) (*model.User, error) {
	u := model.User{}
	query := `SELECT * FROM users WHERE username = $1`
	if err := repo.db.Get(&u, query, username); err != nil {
		return nil, err
	}
	return &u, nil
}

func (repo *UserRepo) Update(id uuid.UUID, u *model.User) error {
	query := `
		UPDATE users
		SET updated_at = $2, email = $3, name = $4
		WHERE id = $1
	`
	_, err := repo.db.Exec(query, id, time.Now().UTC(), u.Email, u.UserName)
	return err
}

func (repo *UserRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
