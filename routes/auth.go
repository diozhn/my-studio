package routes

import (
	"my-studio/database"
	"my-studio/models"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/refresh-token", RefreshToken)
}

func Register(c *fiber.Ctx) error {
	type Input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error encrypting password"})
	}

	user := models.User{
		Username: input.Username,
		Password: string(hashed),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error creating user"})
	}

	return c.JSON(fiber.Map{"message": "User creating success"})
}

func Login(c *fiber.Ctx) error {
	type Input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User or password invalid"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "User or password invalid"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Token generate error"})
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)

	user.RefreshToken = refreshTokenString
	database.DB.Save(&user)

	return c.JSON(fiber.Map{
		"token":         tokenString,
		"refresh_token": refreshTokenString,
	})
}

func RequireAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return c.Status(401).JSON(fiber.Map{"error": "Token not provided"})
	}

	tokenStr := strings.TrimPrefix(auth, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Token"})
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user_id", claims["user_id"].(float64))
		return c.Next()
	}

	return c.Status(401).JSON(fiber.Map{"error": "Invalid Token"})
}

func RefreshToken(c *fiber.Ctx) error {
	type Body struct {
		RefreshToken string `json:"refresh_token"`
	}

	var body Body

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	token, err := jwt.Parse(body.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	userID := uint(claims["user_id"].(float64))

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil || user.RefreshToken != body.RefreshToken {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid refresh token"})
	}

	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	newTokenString, _ := newToken.SignedString(jwtSecret)

	return c.JSON(fiber.Map{"token": newTokenString})
}
