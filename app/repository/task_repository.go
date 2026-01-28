package repository

import (
	"fmt"
	"time"

	"github.com/byeblogs/go-boilerplate/app/model"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/google/uuid"
)

type TaskRepo struct {
	db *database.DB
}

func NewTaskRepo(db *database.DB) TaskRepository {
	return &TaskRepo{db}
}

func (repo *TaskRepo) Create(t *model.Task) error {
	query := `
		INSERT INTO tasks (id, project_id, title, status, due_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	now := time.Now().UTC()
	_, err := repo.db.Exec(query, t.ID, t.ProjectID, t.Title, t.Status, t.DueAt, now, now)
	return err
}

func (repo *TaskRepo) All(limit int, offset uint) ([]*model.Task, error) {
	var out []*model.Task
	query := `SELECT * FROM tasks ORDER BY created_at DESC`
	var err error

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT $1 OFFSET $2", query)
		err = repo.db.Select(&out, query, limit, offset)
	} else {
		err = repo.db.Select(&out, query)
	}
	return out, err
}

func (repo *TaskRepo) AllByProject(projectID uuid.UUID, limit int, offset uint) ([]*model.Task, error) {
	var out []*model.Task
	query := `SELECT * FROM tasks WHERE project_id = $1 ORDER BY created_at DESC`
	var err error

	if limit > 0 {
		query = fmt.Sprintf("%s LIMIT $2 OFFSET $3", query)
		err = repo.db.Select(&out, query, projectID, limit, offset)
	} else {
		err = repo.db.Select(&out, query, projectID)
	}
	return out, err
}

func (repo *TaskRepo) Get(id uuid.UUID) (*model.Task, error) {
	t := model.Task{}
	query := `SELECT * FROM tasks WHERE id = $1`
	if err := repo.db.Get(&t, query, id); err != nil {
		return nil, err
	}
	return &t, nil
}

func (repo *TaskRepo) Update(id uuid.UUID, t *model.Task) error {
	query := `
		UPDATE tasks
		SET updated_at = $2, title = $3, status = $4, due_at = $5
		WHERE id = $1
	`
	_, err := repo.db.Exec(query, id, time.Now().UTC(), t.Title, t.Status, t.DueAt)
	return err
}

func (repo *TaskRepo) Delete(id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := repo.db.Exec(query, id)
	return err
}
