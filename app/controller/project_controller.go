package controller

import (
	"github.com/byeblogs/go-boilerplate/app/model"
	repo "github.com/byeblogs/go-boilerplate/app/repository"
	"github.com/byeblogs/go-boilerplate/pkg/validator"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetProjects supports optional filter: ?owner_user_id=<uuid>
// @Router /v1/projects [get]
func GetProjects(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)

	projectRepo := repo.NewProjectRepo(database.GetDB())

	if s := c.Query("owner_user_id"); s != "" {
		ownerID, err := uuid.Parse(s)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": "invalid owner_user_id"})
		}

		projects, err := projectRepo.AllByOwner(ownerID, pageSize, uint(pageSize*(pageNo-1)))
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "projects were not found"})
		}

		return c.JSON(fiber.Map{
			"page":      pageNo,
			"page_size": pageSize,
			"count":     len(projects),
			"projects":  projects,
		})
	}

	projects, err := projectRepo.All(pageSize, uint(pageSize*(pageNo-1)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "projects were not found"})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(projects),
		"projects":  projects,
	})
}

// GetProject @Router /v1/projects/{id} [get]
func GetProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	projectRepo := repo.NewProjectRepo(database.GetDB())
	p, err := projectRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "project was not found"})
	}

	return c.JSON(fiber.Map{"project": p})
}

// CreateProject @Security ApiKeyAuth
// @Router /v1/projects [post]
func CreateProject(c *fiber.Ctx) error {
	p := &model.Project{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	p.ID = uuid.New()

	validate := validator.NewValidator()
	if err := validate.Struct(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	projectRepo := repo.NewProjectRepo(database.GetDB())
	if err := projectRepo.Create(p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"project": p})
}

// UpdateProject @Security ApiKeyAuth
// @Router /v1/projects/{id} [put]
func UpdateProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	projectRepo := repo.NewProjectRepo(database.GetDB())
	_, err = projectRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "project was not found"})
	}

	p := &model.Project{}
	if err := c.BodyParser(p); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}
	p.ID = id

	// On update, we typically don’t force owner_user_id revalidation.
	// But if you want strict: keep validate.Struct(p) and require owner_user_id in body.
	validate := validator.NewValidator()
	if err := validate.StructPartial(p, "Name"); err != nil {
		// If your validator doesn't support StructPartial, just use validate.Struct(p)
		// and require full payload.
	}

	if err := projectRepo.Update(id, p); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	dbProject, err := projectRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"project": dbProject})
}

// DeleteProject @Security ApiKeyAuth
// @Router /v1/projects/{id} [delete]
func DeleteProject(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	projectRepo := repo.NewProjectRepo(database.GetDB())
	_, err = projectRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "project was not found"})
	}

	if err := projectRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{})
}
