package controller

import (
	"github.com/byeblogs/go-boilerplate/app/model"
	repo "github.com/byeblogs/go-boilerplate/app/repository"
	"github.com/byeblogs/go-boilerplate/pkg/validator"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUsers @Router /v1/users [get]
func GetUsers(c *fiber.Ctx) error {
	pageNo, pageSize := GetPagination(c)

	userRepo := repo.NewUserRepo(database.GetDB())
	users, err := userRepo.All(pageSize, uint(pageSize*(pageNo-1)))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "users were not found"})
	}

	return c.JSON(fiber.Map{
		"page":      pageNo,
		"page_size": pageSize,
		"count":     len(users),
		"users":     users,
	})
}

// GetUser @Router /v1/users/{id} [get]
func GetUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	u, err := userRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "user was not found"})
	}

	return c.JSON(fiber.Map{"user": u})
}

// CreateUser @Security ApiKeyAuth
// @Router /v1/users [post]
func CreateUser(c *fiber.Ctx) error {
	u := &model.User{}
	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	u.ID = uuid.New()

	validate := validator.NewValidator()
	if err := validate.Struct(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	if err := userRepo.Create(u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"user": u})
}

// UpdateUser @Security ApiKeyAuth
// @Router /v1/users/{id} [put]
func UpdateUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	_, err = userRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "user was not found"})
	}

	u := &model.User{}
	if err := c.BodyParser(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}
	u.ID = id

	validate := validator.NewValidator()
	if err := validate.Struct(u); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg":    "invalid input found",
			"errors": validator.ValidatorErrors(err),
		})
	}

	if err := userRepo.Update(id, u); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	dbUser, err := userRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{"user": dbUser})
}

// DeleteUser @Security ApiKeyAuth
// @Router /v1/users/{id} [delete]
func DeleteUser(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"msg": err.Error()})
	}

	userRepo := repo.NewUserRepo(database.GetDB())
	_, err = userRepo.Get(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"msg": "user was not found"})
	}

	if err := userRepo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"msg": err.Error()})
	}

	return c.JSON(fiber.Map{})
}
