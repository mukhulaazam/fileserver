package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func main() {
	app := fiber.New()

	app.Post("/api/upload", func(c *fiber.Ctx) error {
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}
		
		uniqueID := uuid.New().String()
		extension := filepath.Ext(file.Filename)
		uniqueFilename := fmt.Sprintf("%s%s", uniqueID, extension)

		currentTime := time.Now()
		folderPath := fmt.Sprintf("uploads/%02d-%d-%d", currentTime.Day(), currentTime.Month(), currentTime.Year())

		if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create folder",
				"error":   err.Error(),
			})
		}

		// Save the file to the created folder with the unique filename
		savePath := filepath.Join(folderPath, uniqueFilename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "File upload failed",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message":  "File uploaded successfully",
			"filepath": "/" + folderPath + uniqueFilename,
		})
	})
	app.Listen(":3000")
}
