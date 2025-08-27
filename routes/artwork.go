package routes

import (
	"fmt"
	"my-studio/database"
	"my-studio/middlewares"
	"my-studio/models"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func RegisterArtworkRoutes(app *fiber.App) {
	app.Post("/artworks", RequireAuth, middlewares.IsArtworkOwner, CreateArtwork)
	app.Post("/artworks/:id/like", LikeArtwork)
	app.Get("/gallery", GetGallery)
	app.Get("/top-artworks", GetTopArtworks)
	app.Get("/artworks", GetArtworks)
	app.Get("/artworks/:id", GetArtworksById)
	app.Get("/artworks/filter", GetFilteredArtworks)
	app.Patch("/artworks/:id", RequireAuth, middlewares.IsArtworkOwner, EditArtwork)
	app.Delete("/artworks/:id", RequireAuth, middlewares.IsArtworkOwner, DeleteArtwork)
}

func GetArtworks(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	title := c.Query("title")
	author := c.Query("author")
	from := c.Query("from")
	to := c.Query("to")
	sort := c.Query("sort", "created_desc")

	dbQuery := database.DB.Model(&models.Artwork{})

	if title != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+title+"%")
	}
	if author != "" {
		dbQuery = dbQuery.Where("user_id = ?", author)
	}
	if from != "" && to != "" {
		dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", from, to)
	}

	switch sort {
	case "likes_desc":
		dbQuery = dbQuery.Order("likes desc")
	case "likes_asc":
		dbQuery = dbQuery.Order("likes asc")
	case "created_asc":
		dbQuery = dbQuery.Order("created_at asc")
	default:
		dbQuery = dbQuery.Order("created_at desc")
	}

	var artworks []models.Artwork
	if err := dbQuery.Limit(limit).Offset(offset).Find(&artworks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error fetching artworks"})
	}

	return c.JSON(fiber.Map{
		"page":    page,
		"limit":   limit,
		"results": artworks,
		"count":   len(artworks),
	})
}

func GetGallery(c *fiber.Ctx) error {
	var artworks []models.Artwork

	if err := database.DB.Order("created_at DESC").Find(&artworks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve artworks: " + err.Error(),
		})
	}

	html := "<html><body><h1>Galeria de Artes</h1><div style='display: flex; flex-wrap: wrap;'>"

	for _, art := range artworks {
		html += fmt.Sprintf(`
					<div style='margin:10px; text-align:center'>
				<img src='%s' style='max-width:200px; max-height:200px; display:block;' />
				<h3>%s</h3>
				<p>%s</p>
				<p>Likes: %d</p>
			</div>
		`, art.ImageURL, art.Title, art.Caption, art.Likes)
	}

	html += "</div>"

	return c.Type("html").SendString(html)
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

	userID := c.Locals("user_id").(uint)
	art := models.Artwork{
		Title:     title,
		Caption:   caption,
		ImageURL:  "/" + filePath,
		UserID:    userID,
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

func LikeArtwork(c *fiber.Ctx) error {
	id := c.Params("id")
	var artwork models.Artwork

	if err := database.DB.First(&artwork, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Artwork not found",
		})
	}

	artwork.Likes = artwork.Likes + 1

	if err := database.DB.Save(&artwork).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to like artwork: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Artwork liked successfully",
		"likes":   artwork.Likes,
	})
}

func GetTopArtworks(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	var artworks []models.Artwork
	if err := database.DB.Order("likes DESC").Limit(limit).Find(&artworks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve top artworks: " + err.Error(),
		})
	}

	return c.JSON(artworks)
}

func GetFilteredArtworks(c *fiber.Ctx) error {
	title := c.Query("title")
	author := c.Query("author")
	from := c.Query("from")
	to := c.Query("to")

	dbQuery := database.DB.Model(&models.Artwork{})

	if title != "" {
		dbQuery = dbQuery.Where("title LIKE ?", "%"+title+"%")
	}

	if author != "" {
		dbQuery = dbQuery.Where("user_id = ?", author)
	}

	if from != "" && to != "" {
		dbQuery = dbQuery.Where("created_at BETWEEN ? AND ?", from, to)
	}

	var artworks []models.Artwork
	if err := dbQuery.Order("created_at DESC").Find(&artworks).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve artworks: " + err.Error(),
		})
	}

	return c.JSON(artworks)
}
