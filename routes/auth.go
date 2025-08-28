package routes

import (
	"my-studio/database"
	"my-studio/models"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/shareed2k/goth_fiber"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func RegisterAuthRoutes(app *fiber.App) {
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/refresh-token", RefreshToken)
	app.Get("/auth/:provider", SocialAuthBegin)
	app.Get("/auth/:provider/callback", SocialAuthCallback)
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

// Inicia o fluxo OAuth
func SocialAuthBegin(c *fiber.Ctx) error {
	err := goth_fiber.BeginAuthHandler(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return nil
}

// Callback do OAuth
func SocialAuthCallback(c *fiber.Ctx) error {
	provider := c.Params("provider")
	c.Context().SetUserValue("provider", provider)
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	var dbUser models.User
	// Procura usuário pelo ID da rede social
	switch provider {
	case "google":
		database.DB.Where("google_id = ?", user.UserID).First(&dbUser)
	case "instagram":
		database.DB.Where("instagram_id = ?", user.UserID).First(&dbUser)
	case "twitter":
		database.DB.Where("twitter_id = ?", user.UserID).First(&dbUser)
	}

	type UserResponse struct {
		ID          uint   `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		GoogleID    string `json:"google_id,omitempty"`
		InstagramID string `json:"instagram_id,omitempty"`
		TwitterID   string `json:"twitter_id,omitempty"`
	}

	// Se não existir, cria novo usuário
	if dbUser.ID == 0 {
		dbUser = models.User{
			Username:    user.Name,
			Email:       user.Email,
			GoogleID:    ifThenElse(provider == "google", user.UserID, ""),
			InstagramID: ifThenElse(provider == "instagram", user.UserID, ""),
			TwitterID:   ifThenElse(provider == "twitter", user.UserID, ""),
		}
		database.DB.Create(&dbUser)
	}

	userResponse := UserResponse{
		ID:          dbUser.ID,
		Username:    dbUser.Username,
		Email:       dbUser.Email,
		GoogleID:    dbUser.GoogleID,
		InstagramID: dbUser.InstagramID,
		TwitterID:   dbUser.TwitterID,
	}

	// Gera JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": dbUser.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	return c.JSON(fiber.Map{
		"token": tokenString,
		"user":  userResponse,
	})
}

// Função auxiliar para atribuição condicional
func ifThenElse(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}
