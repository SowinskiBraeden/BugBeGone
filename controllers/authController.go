package controllers

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/SowinskiBraeden/BugBeGone/database"
	"github.com/SowinskiBraeden/BugBeGone/models"
	"github.com/SowinskiBraeden/gfbmb/messageBox"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	firstname := c.FormValue("firstname")
	lastname := c.FormValue("lastname")
	username := c.FormValue("username")
	email := c.FormValue("email")
	password := c.FormValue("password")

	if email == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": messageBox.NewWarningBox("Missing Email"),
			"username": username,
			"email":    email,
		})
	}
	if password == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": messageBox.NewWarningBox("Please enter your password"),
			"username": username,
			"email":    email,
		})
	}
	if username == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": messageBox.NewWarningBox("Missing Username"),
			"username": username,
			"email":    email,
		})
	}

	// Check if email or username is previously registered
	count, err := userCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			bson.M{"email": email},
			bson.M{"username": username},
		}},
	)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
			"errorMsg": messageBox.NewWarningBox("Failed to search database"),
			"username": username,
			"email":    email,
		})
	}
	if count > 0 {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
			"errorMsg": messageBox.NewWarningBox("An account already exists with that email or username"),
			"username": username,
			"email":    email,
		})
	}

	var user models.User
	if user.CheckPasswordStrength(password) {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": messageBox.NewWarningBox("Your password must contain at least 1 lowercase, 1 uppercase & 1 special character"),
			"username": username,
			"email":    email,
		})
	}

	user.Username = username
	user.Firstname = firstname
	user.Lastname = lastname
	user.Email = email
	user.Password = user.HashPassword(password)
	user.TempPassword = false
	user.ID = primitive.NewObjectID()
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be inserted",
			"error":   insertErr,
		})
	}

	defer cancel()

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{})
}
