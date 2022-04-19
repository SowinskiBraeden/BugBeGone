package controllers

import (
	"context"

	"github.com/SowinskiBraeden/BugBeGone/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// Render Index
func MainPage(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	// This returns not authorized for both admin and student
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).Render("index", fiber.Map{
			"loggedIn": false,
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	uid := claims.Issuer

	var user models.User
	findErr := userCollection.FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&user)
	if findErr != nil {
		return c.Status(fiber.StatusUnauthorized).Render("index", fiber.Map{
			"loggedIn": false,
		})
	}

	return c.Status(fiber.StatusAccepted).Render("index", fiber.Map{
		"loggedIn": true,
	})
}

func RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", fiber.Map{
		"errorMsg": "",
	})
}

func LoginPage(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	// This returns not authorized
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	uid := claims.Issuer

	var user models.User
	findErr := userCollection.FindOne(context.TODO(), bson.M{"uid": uid}).Decode(&user)
	if findErr != nil {
		return c.Status(fiber.StatusUnauthorized).Render("login", fiber.Map{
			"msg":      "",
			"errorMsg": "",
		})
	}

	return c.Status(fiber.StatusAccepted).Render("dashboard", fiber.Map{
		"msg":       "",
		"errorMsg":  "",
		"username":  user.Username,
		"email":     user.Email,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
	})
}
