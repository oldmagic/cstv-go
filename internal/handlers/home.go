package handlers

import "github.com/gofiber/fiber/v2"

// HomeHandler serves the home page
func HomeHandler(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{"message": "Welcome to GOTV-Plus"})
}
