package controllers

import (
	"github.com/SowinskiBraeden/gfbmb/messageBox"
	"github.com/gofiber/fiber/v2"
)

// Render Index
func MainPage(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{})
}

func RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"errorMsg": messageBox.EmptyMessageBox(),
	})
}

func LoginPage(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"errorMsg": messageBox.EmptyMessageBox(),
	})
}
