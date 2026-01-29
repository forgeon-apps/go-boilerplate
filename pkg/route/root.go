package route

import (
	"strings"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/byeblogs/go-boilerplate/platform/database"
	"github.com/gofiber/fiber/v2"
)

func GeneralRoute(a *fiber.App) {

	a.Get("/", func(c *fiber.Ctx) error {
		accept := c.Get("Accept")
		if strings.Contains(accept, "application/json") {
			return c.JSON(fiber.Map{
				"msg":    "Welcome to Fiber Go API!",
				"docs":   "/swagger/index.html",
				"status": "/h34l7h",
				"ui":     "/api/v1/ui",
			})
		}
		return c.Redirect("/api/v1/ui", fiber.StatusFound)
	})

	a.Get("/healthz", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"msg":    "Welcome to Fiber Go API!",
			"docs":   "/swagger/index.html",
			"status": "/h34l7h",
		})
	})

	a.Get("/h34l7h", func(c *fiber.Ctx) error {
		err := database.GetDB().Ping()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"msg": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"msg":       "Health Check",
			"db_online": true,
		})
	})
}

func SwaggerRoute(a *fiber.App) {
	// Create route group.
	route := a.Group("/swagger")
	route.Get("*", swagger.Handler)
}

func NotFoundRoute(a *fiber.App) {
	a.Use(
		func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": "sorry, endpoint is not found",
			})
		},
	)
}
