package controller

import (
	"strings"

	"github.com/byeblogs/go-boilerplate/app/model"
	repo "github.com/byeblogs/go-boilerplate/app/repository"
	"github.com/byeblogs/go-boilerplate/pkg/validator"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetTasks supports optional filter: ?project_id=<uuid>
// @Router /v1/tasks [get]
func GetTasks(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)

	taskRepo := repo.NewTaskRepo(database.GetDB())

	if s := c.Query("project_id"); s != "" {
		projectID, err := uuid.Parse(s)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "invalid project_id"})
		}

		tasks, err := taskRepo.AllByProject(projectID, pageSize, uint(pageSize*(pageNo-1)))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "tasks were not found"})
		}

		return c.JSON(fiber.Map{
			"page":      pageNo,
			"page_size": pageSize,
			"count":     len(tasks),
			"tasks":     tasks,
		})
	}

	tasks, err := taskRepo.All(pageSize, uint(pageSize*(pageNo-1)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "tasks were not found"})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(tasks),
		"tasks":     tasks,
	})
}

// GetTask @Router /v1/tasks/{id} [get]
func GetTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	taskRepo := repo.NewTaskRepo(database.GetDB())
	t, err := taskRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "task was not found"})
	}

	return c.JSON(fiber.Map{"task": t})
}

// CreateTask @Security ApiKeyAuth
// @Router /v1/tasks [post]
func CreateTask(c *fiber.Ctx) error {
	t := &model.Task{}
	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	t.ID = uuid.New()
	if strings.TrimSpace(t.Status) == "" {
		t.Status = "todo"
	}

	validate := validator.NewValidator()
	if err := validate.Struct(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	taskRepo := repo.NewTaskRepo(database.GetDB())
	if err := taskRepo.Create(t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"task": t})
}

// UpdateTask @Security ApiKeyAuth
// @Router /v1/tasks/{id} [put]
func UpdateTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	taskRepo := repo.NewTaskRepo(database.GetDB())
	_, err = taskRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "task was not found"})
	}

	t := &model.Task{}
	if err := c.BodyParser(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}
	t.ID = id

	if strings.TrimSpace(t.Status) == "" {
		t.Status = "todo"
	}

	validate := validator.NewValidator()
	if err := validate.Struct(t); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	if err := taskRepo.Update(id, t); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	dbTask, err := taskRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"task": dbTask})
}

// DeleteTask @Security ApiKeyAuth
// @Router /v1/tasks/{id} [delete]
func DeleteTask(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	taskRepo := repo.NewTaskRepo(database.GetDB())
	_, err = taskRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "task was not found"})
	}

	if err := taskRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{})
}
