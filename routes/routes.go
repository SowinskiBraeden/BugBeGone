package routes

import (
	"github.com/SowinskiBraeden/BugBeGone/controllers"
	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {

	// Page handlers
	app.Get("/", controllers.MainPage)
	app.Get("/register", controllers.RegisterPage)
	app.Get("/login", controllers.LoginPage)

	// Authentication handlers
	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.Login)

	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
}
