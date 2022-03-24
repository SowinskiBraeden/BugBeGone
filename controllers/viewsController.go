package controllers

import "github.com/gofiber/fiber/v2"

// Render Index
func MainPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"errorMsg": "<div></div>",
	})
}

func LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"errorMsg": `<div class="alert alert-warning alert-dismissible fade show" role="alert">
		test error
		<button type="button" class="close" data-dismiss="alert" aria-label="Close">
			  <span aria-hidden="true">&times;</span>
		</button>
	</div>`,
	})
}
