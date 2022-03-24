package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/SowinskiBraeden/BugBeGone/database"
	"github.com/SowinskiBraeden/BugBeGone/models"
	"github.com/SowinskiBraeden/BugBeGone/render"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var imageCollection *mongo.Collection = database.OpenCollection(database.Client, "images")

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	email := c.FormValue("email")
	pass1 := c.FormValue("password")
	pass2 := c.FormValue("password2")
	username := c.FormValue("username")
	accountType := c.FormValue("accountType")
	file, _ := c.FormFile("image")

	if email == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("Missing Email"),
			"username": username,
			"email":    email,
		})
	}
	if pass1 == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("Missing Password"),
			"username": username,
			"email":    email,
		})
	}
	if pass2 == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("Please retype your password"),
			"username": username,
			"email":    email,
		})
	}
	if username == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("Missing Username"),
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
			"errorMsg": render.NewWarningBox("Failed to search database"),
			"username": username,
			"email":    email,
		})
	}
	if count > 0 {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("An account already exists with that email or username"),
			"username": username,
			"email":    email,
		})
	}

	if pass1 != pass2 {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("The passwords don't match"),
			"username": username,
			"email":    email,
		})
	}

	var user models.User
	if user.CheckPasswordStrength(pass1) {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("Your password must contain at least 1 lowercase, 1 uppercase & 1 special character"),
			"username": username,
			"email":    email,
		})
	}

	accountType = strings.ToUpper(accountType)
	if accountType != "T" && accountType != "R" && accountType != "A" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": render.NewWarningBox("There appears to be an error with your account type, please try again"),
			"username": username,
			"email":    email,
		})
	}

	user.Username = username
	user.Email = email
	user.Password = user.HashPassword(pass1)
	user.TempPassword = false
	user.AccountType = accountType
	user.ID = primitive.NewObjectID()
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	var photo models.Photo
	if file != nil {
		// Save image to local
		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := strings.Split(file.Filename, ".")[1]
		image := fmt.Sprintf("%s.%s", filename, fileExt)
		err := c.SaveFile(file, fmt.Sprintf("./database/images/%s", image))
		if err != nil {
			cancel()
			return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
				"errorMsg": "the image could not be saved",
				"username": username,
				"email":    email,
			})
		}

		// Read the entire file into a byte slice
		bytes, err := ioutil.ReadFile(fmt.Sprintf("./database/images/%s", image))
		if err != nil {
			cancel()
			return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
				"errorMsg": "the image could not be read",
				"username": username,
				"email":    email,
			})
		}

		var base64Encoding string = toBase64(bytes)

		photo.Name = uuid.New().String()
		photo.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		photo.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		photo.ID = primitive.NewObjectID()

		photo.Base64 = base64Encoding

		// Remove local image
		os.Remove(fmt.Sprintf("./database/images/%s", image))
	} else {
		photo.Name = uuid.New().String()
		photo.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		photo.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		photo.ID = primitive.NewObjectID()

		defaultImage, _ := os.ReadFile("./database/defaultImage.txt")
		photo.Base64 = string(defaultImage)

	}
	user.PhotoName = photo.Name

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student could not be inserted",
			"error":   insertErr,
		})
	}

	_, insertErr = imageCollection.InsertOne(ctx, photo)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "the student default photo could not be inserted",
			"error":   insertErr,
		})
	}
	defer cancel()

	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{})
}

func Login(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{})
}
