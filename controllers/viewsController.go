package controllers

import (
	"github.com/gofiber/fiber/v2"
)

// Render Index
func MainPage(c *fiber.Ctx) error {
	authorized, _ := IsAuthorized(c)
	// If the user is already authorized, redirect to dashboard
	if authorized {
		return c.Status(fiber.StatusAccepted).Render("index", fiber.Map{
			"loggedIn": true,
		})
	}
	return c.Status(fiber.StatusAccepted).Render("index", fiber.Map{
		"loggedIn": false,
	})
}

func RegisterPage(c *fiber.Ctx) error {
	authorized, _ := IsAuthorized(c)
	// If the user is already authorized, redirect to dashboard
	if authorized {
		return c.Redirect("/dashboard")
	}
	return c.Render("register", fiber.Map{
		"errorMsg": "",
	})
}

func LoginPage(c *fiber.Ctx) error {
	authorized, user := IsAuthorized(c)
	// If the user is already authorized, redirect to dashboard
	if authorized {
		return c.Status(fiber.StatusAccepted).Redirect("/dashboard")
	}
	return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
		"errorMsg": "",
		"msg":      "",
		"user":     user,
	})
}

func DashboardPage(c *fiber.Ctx) error {
	authorized, user := IsAuthorized(c)
	if !authorized {
		return c.Status(fiber.StatusUnauthorized).Redirect("/login")
	}
	return c.Status(fiber.StatusAccepted).Render("dashboard", fiber.Map{
		"errorMsg": "",
		"msg":      "",
		"user":     user,
	})
}

func ProfilePage(c *fiber.Ctx) error {
	authorized, user := IsAuthorized(c)
	if !authorized {
		return c.Status(fiber.StatusUnauthorized).Redirect("/login")
	}
	return c.Status(fiber.StatusAccepted).Render("profile", fiber.Map{
		"errorMsg": "",
		"msg":      "",
		"user":     user,
	})
}
