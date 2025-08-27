package routes

import (
	"my-studio/database"
	"my-studio/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/shareed2k/goth_fiber"
	"golang.org/x/crypto/bcrypt"
)

func UserRoutes(app *fiber.App) {
	app.Get("/users/:id", RequireAuth, GetUserProfile)
	app.Patch("/users/:id", RequireAuth, UpdateUserProfile)
	app.Get("/users/:id/artworks", GetUserArtworks)
	app.Post("/users/:provider/link", RequireAuth, LinkSocialAccount)
}

func GetUserProfile(c *fiber.Ctx) error {
	id := c.Params("id")

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":         user.ID,
		"username":   user.Username,
		"created_at": user.CreatedAt,
	})
}

func UpdateUserProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	paramID, _ := strconv.Atoi(c.Params("id"))

	if uint(paramID) != userID {
		return c.Status(403).JSON(fiber.Map{"error": "You can only update your own profile"})
	}

	type Input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var input Input
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if input.Username != "" {
		user.Username = input.Username
	}

	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		user.Password = string(hashed)
	}

	database.DB.Save(&user)
	return c.JSON(fiber.Map{"message": "Profile updated successfully"})
}

func GetUserArtworks(c *fiber.Ctx) error {
	userID, _ := strconv.Atoi(c.Params("id"))

	var artworks []models.Artwork
	if err := database.DB.Where("user_id = ?", userID).Find(&artworks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching artworks"})
	}

	return c.JSON(artworks)
}

func LinkSocialAccount(c *fiber.Ctx) error {
	provider := c.Params("provider")
	c.Context().SetUserValue("provider", provider)
	user, err := goth_fiber.CompleteUserAuth(c)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	userID := c.Locals("user_id").(uint)
	var dbUser models.User
	if err := database.DB.First(&dbUser, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	switch provider {
	case "google":
		dbUser.GoogleID = user.UserID
	case "instagram":
		dbUser.InstagramID = user.UserID
	case "twitter":
		dbUser.TwitterID = user.UserID
	}
	database.DB.Save(&dbUser)

	return c.JSON(fiber.Map{"message": "Conta vinculada com sucesso!", "user": dbUser})
}
