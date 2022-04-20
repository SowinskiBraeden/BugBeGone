package controllers

import (
	"context"
	"os"
	"time"

	"github.com/SowinskiBraeden/BugBeGone/database"
	"github.com/SowinskiBraeden/BugBeGone/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"

	"github.com/dgrijalva/jwt-go"

	"github.com/gofiber/fiber/v2"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "users")
var SecretKey string

func Init() {
	godotenv.Load(".env")
	SecretKey = os.Getenv("secret")
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
			"errorMsg": "Missing Email",
		})
	}
	if password == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": "Please enter your password",
		})
	}
	if username == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": "Missing Username",
		})
	}

	// Check if email or username is previously registered
	count, err := userCollection.CountDocuments(ctx, bson.M{
		"$or": []bson.M{
			{"email": email},
			{"username": username},
		}},
	)
	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
			"errorMsg": "Failed to search database",
		})
	}
	if count > 0 {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
			"errorMsg": "An account already exists with that email or username",
		})
	}

	var user models.User
	if user.CheckPasswordStrength(password) {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("register", fiber.Map{
			"errorMsg": "Your password must contain at least 1 lowercase, 1 uppercase & 1 special character",
			"username": username,
			"email":    email,
		})
	}

	user.UID = uuid.New().String()
	user.Username = username
	user.Firstname = firstname
	user.Lastname = lastname
	user.Email = email
	user.Password = user.HashPassword(password)
	user.TempPassword = false
	user.Attempts = 0
	user.Disabled = false
	user.ID = primitive.NewObjectID()
	user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	_, insertErr := userCollection.InsertOne(ctx, user)
	if insertErr != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("register", fiber.Map{
			"errorMsg": "the user could not be inserted",
		})
	}

	defer cancel()

	return c.Render("login", fiber.Map{
		"msg":      "Successfully registered an account",
		"errorMsg": "",
	})
}

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" && password == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "Username and password can't be blank",
		})
	}
	if username == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "Username can't be blank",
		})
	}
	if password == "" {
		cancel()
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "Password can't be blank",
		})
	}

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{
		"$or": []bson.M{
			{"email": username},
			{"username": username},
		},
	}).Decode(&user)
	defer cancel()

	if err != nil {
		cancel()
		return c.Status(fiber.StatusInternalServerError).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "user not found",
		})
	}

	var localAccountDisabled = false
	if user.Attempts >= 5 {
		localAccountDisabled = true // Catches newly disbaled account before student obj is updated
		update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		update := bson.M{
			"$set": bson.M{
				"disabled":   true,
				"attempts":   0,
				"updated_at": update_time,
			},
		}

		_, updateErr := userCollection.UpdateOne(
			ctx,
			bson.M{
				"$or": []bson.M{
					{"email": username},
					{"username": username},
				},
			},
			update,
		)
		if updateErr != nil {
			cancel()
			return c.Status(fiber.StatusInternalServerError).Render("login", fiber.Map{
				"msg":      "",
				"errorMsg": "the user could not be updated",
			})
		}
	}

	if localAccountDisabled || user.Disabled {
		cancel()
		return c.Status(fiber.StatusForbidden).Render("login", fiber.Map{
			"msg":      "Account is Disabled, contact support",
			"errorMsg": "Account is Disabled, contact support",
		})
	}

	var verified bool = user.ComparePasswords(password)
	if verified == false {
		update_time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		update := bson.M{
			"$set": bson.M{
				"Attempts":   user.Attempts + 1,
				"updated_at": update_time,
			},
		}

		_, updateErr := userCollection.UpdateOne(
			ctx,
			bson.M{"username": username},
			update,
		)
		cancel()
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).Render("login", fiber.Map{
				"msg":      "",
				"errorMsg": "the student could not be updated",
			})
		}
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "incorrect password",
		})
	}
	defer cancel()

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.UID,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // 1 Day
	})
	token, err := claims.SignedString([]byte(SecretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "could not log in",
		})
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusNotImplemented).Redirect("/dashboard")
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.Status(fiber.StatusOK).Redirect("/")
}
