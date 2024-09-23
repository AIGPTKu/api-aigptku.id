package rest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (r *restHandler) uploadFileTemp(c *fiber.Ctx) (err error) {

	// Get the file from the form
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusBadRequest).SendString("Cannot get the file")
	}

	// Generate a UUID and get the file extension
	filename := uuid.New().String() + filepath.Ext(file.Filename)

	// Create a new file on the server
	dst, err := os.Create("assets/temp/" + filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot create file")
	}
	defer dst.Close()

	// Save the file
	if err := c.SaveFile(file, dst.Name()); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Cannot save file")
	}

	return c.JSON(fiber.Map{
		"code": 200,
		"status": "OK",
		"data": fiber.Map{
			"url": viper.GetString("DOMAIN") + "/files/temp/" + filename,
		},
	})
}