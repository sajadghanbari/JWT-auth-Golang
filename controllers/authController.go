package controllers

import (
	"JWT-Authentication-go/database"
	"JWT-Authentication-go/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

 func Hello(c *fiber.Ctx) error {
     return c.SendString("Hello world!!")
 }

func Register (c *fiber.Ctx) error{
	fmt.Println("Received a registration request")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
    var existingUser models.User
    if err := database.DB.Where("email = ?", data["email"]).First(&existingUser).Error; err == nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Email already exists",
        })
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to hash password",
   		})
	}
	user := models.User{
		Name: data["name"],
		Email: data["email"],
		Password: string(hashedPassword),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}