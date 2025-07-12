package routes

import (
	"fmt"
	"my-studio/database"
	"my-studio/models"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterArtworkRoutes(app *fiber.App) {
	app.Get("/artworks", GetArtWorks)
	app.Post("/artworks", CreateArtwork)
	app.Get("/artworks/:id", GetArtworksById)
	app.Delete("/artworks/:id", DeleteArtwork)
	app.Patch("/artworks/:id", EditArtwork)
}

func GetArtWorks(c *fiber.Ctx) error {
	var artworks []models.Artwork

	database.DB.Order("created_at DESC").Find(&artworks)

	return c.JSON(artworks)
}

func CreateArtwork(c *fiber.Ctx) error {
	title := c.FormValue("title")
	caption := c.FormValue("caption")

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Bad request: " + err.Error(),
		})
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join("uploads", filename)

	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file: " + err.Error(),
		})
	}

	art := models.Artwork {
		Title: 	 title,
		Caption: caption,
		ImageURL: "/" + filePath,
		CreatedAt: time.Now(),
	}

	result := database.DB.Create(&art)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create artwork: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(art)
}

func GetArtworksById(c *fiber.Ctx) error {
	id := c.Params("id")
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Artwork not found",
		})
	}

	return c.JSON(artwork)
}

func DeleteArtwork(c *fiber.Ctx) error {
	id := c.Params("id")
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Artwork not found",
		})
	}

	if artwork.ImageURL != "" {
		err := os.Remove("." + artwork.ImageURL)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete image file: " + err.Error(),
			})
		}
	}

	if err := database.DB.Delete(&artwork, id).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete artwork: " + err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}


type EditArtworkInput struct {
	Title   string `json:"title"`
	Caption string `json:"caption"`
}

func EditArtwork(c *fiber.Ctx) error {
	id := c.Params("id")
	var artwork models.Artwork
	var input EditArtworkInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON: " + err.Error(),
		})
	}

	if err := database.DB.First(&artwork, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Artwork not found",
		})
	}

	if input.Title != "" {
		artwork.Title = input.Title
	}

	if input.Caption != "" {
		artwork.Caption = input.Caption
	}

	if err := database.DB.Save(&artwork).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update artwork: " + err.Error(),
		})
	}

	return c.JSON(artwork)
}