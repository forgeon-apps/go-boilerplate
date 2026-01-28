package repository

import (
	"fmt"
	"time"

	"github.com/byeblogs/go-boilerplate/app/model"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/google/uuid"
)

type ProjectRepo struct {
	db *database.DB
}

func NewProjectRepo(db *database.DB) ProjectRepository {
	return &ProjectRepo{db}
}

func (repo *ProjectRepo) Create(p *model.Project) error {
	query := `
		INSERT INTO projects (id, owner_user_id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	now := time.Now().UTC()
	_, err := repo.db.Exec(query, p.ID, p.OwnerUserID, p.Name, p.Description, now, now)
	return err
}

func (repo *ProjectRepo) All(limit int, offset uint) ([]*model.Project, error) {
	var out []*model.Project
	query := `SELECT * FROM projects ORDER BY created_at DESC`
	var err error

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
		err = repo.db.Select(&out, query, limit, offset)
	} else {
		err = repo.db.Select(&out, query)
	}
	return out, err
}

func (repo *ProjectRepo) AllByOwner(ownerID uuid.UUID, limit int, offset uint) ([]*model.Project, error) {
	var out []*model.Project
	query := `SELECT * FROM projects WHERE owner_user_id = $1 ORDER BY created_at DESC`
	var err error

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT $2 OFFSET $3", query)
		err = repo.db.Select(&out, query, ownerID, limit, offset)
	} else {
		err = repo.db.Select(&out, query, ownerID)
	}
	return out, err
}

func (repo *ProjectRepo) Get(id uuid.UUID) (*model.Project, error) {
	p := model.Project{}
	query := `SELECT * FROM projects WHERE id = $1`
	if err := repo.db.Get(&p, query, id); err != nil {
		return nil, err
	}
	return &p, nil
}

func (repo *ProjectRepo) Update(id uuid.UUID, p *model.Project) error {
	query := `
		UPDATE projects
		SET updated_at = $2, name = $3, description = $4
		WHERE id = $1
	`
	_, err := repo.db.Exec(query, id, time.Now().UTC(), p.Name, p.Description)
	return err
}

func (repo *ProjectRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
