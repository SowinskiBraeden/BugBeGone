package main

import (
	"os"

	"github.com/SowinskiBraeden/BugBeGone/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html"
	"github.com/joho/godotenv"
)

func main() {
	// Create a new engine
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
