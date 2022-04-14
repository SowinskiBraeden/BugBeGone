package main

import (
	"html/template"
	"os"

	"github.com/SowinskiBraeden/BugBeGone/controllers"
	"github.com/SowinskiBraeden/BugBeGone/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	// Create a new engine
	engine := html.New("./views", ".html")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

	controllers.Init()

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./public")

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
