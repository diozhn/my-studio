package middlewares

import (
	"my-studio/database"
	"my-studio/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func IsArtworkOwner(c *fiber.Ctx) error {
	artworkID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid artwork ID"})
	}

	var artwork models.Artwork
	if err := database.DB.First(&artwork, artworkID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Artwork not found"})
	}

	userID := c.Locals("user_id").(uint)
	if artwork.UserID != userID {
		return c.Status(403).JSON(fiber.Map{"error": "You are not the owner of this artwork"})
	}
	return c.Next()
}
