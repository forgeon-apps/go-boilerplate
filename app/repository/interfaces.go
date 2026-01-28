package repository

import (
	"github.com/byeblogs/go-boilerplate/app/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(u *model.User) error
	All(limit int, offset uint) ([]*model.User, error)
	Get(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(id uuid.UUID, u *model.User) error
	Delete(id uuid.UUID) error
}

type BookRepository interface {
	Create(b *model.Book) error
	All(limit int, offset uint) ([]*model.Book, error)
	Get(ID uuid.UUID) (*model.Book, error)
	Update(ID uuid.UUID, b *model.Book) error
	Delete(ID uuid.UUID) error
}

type ProjectRepository interface {
	Create(p *model.Project) error
	All(limit int, offset uint) ([]*model.Project, error)
	AllByOwner(ownerID uuid.UUID, limit int, offset uint) ([]*model.Project, error)
	Get(id uuid.UUID) (*model.Project, error)
	Update(id uuid.UUID, p *model.Project) error
	Delete(id uuid.UUID) error
}

type TaskRepository interface {
	Create(t *model.Task) error
	All(limit int, offset uint) ([]*model.Task, error)
	AllByProject(projectID uuid.UUID, limit int, offset uint) ([]*model.Task, error)
	Get(id uuid.UUID) (*model.Task, error)
	Update(id uuid.UUID, t *model.Task) error
	Delete(id uuid.UUID) error
}
