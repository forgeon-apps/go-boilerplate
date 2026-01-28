package route

import (
	"github.com/byeblogs/go-boilerplate/app/controller"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public route.
func PublicRoutes(a *fiber.App) {
	// Create route group.
	route := a.Group("/api/v1")

	route.Post("/token/new", controller.GetNewAccessToken)
	route.Get("/books", controller.GetBooks)
	route.Get("/books/:id", controller.GetBook)

	// Users
	route.Get("/users", controller.GetUsers)
	route.Get("/users/:id", controller.GetUser)
	route.Post("/users", controller.CreateUser)
	route.Put("/users/:id", controller.UpdateUser)
	route.Delete("/users/:id", controller.DeleteUser)

	// Projects
	route.Get("/projects", controller.GetProjects)
	route.Get("/projects/:id", controller.GetProject)
	route.Post("/projects", controller.CreateProject)
	route.Put("/projects/:id", controller.UpdateProject)
	route.Delete("/projects/:id", controller.DeleteProject)

	// Tasks
	route.Get("/tasks", controller.GetTasks)
	route.Get("/tasks/:id", controller.GetTask)
	route.Post("/tasks", controller.CreateTask)
	route.Put("/tasks/:id", controller.UpdateTask)
	route.Delete("/tasks/:id", controller.DeleteTask)

}
