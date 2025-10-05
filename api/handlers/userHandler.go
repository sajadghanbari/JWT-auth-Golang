// handlers/item_handler.go
package handlers

import (
	"JWT-Authentication-go/api/dto"
	"JWT-Authentication-go/data/models"
	"JWT-Authentication-go/services"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var validate = validator.New()

func CreateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.CreateUpdateUserRequest
		err := c.BodyParser(&req)
		if err != nil{
			return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}

		err = validate.Struct(req)
		if err != nil{
			return  c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		hashedPassword := services.HashPassword(req.Password)
		newUser := models.User{
			Name: req.Name,
			Email: req.Email,
			Password: hashedPassword,
		}

		err = db.Create(&newUser).Error
		if err != nil {
			    if errors.Is(err, gorm.ErrDuplicatedKey) {
                return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                    "error": "User with this email already exists",
                })
            }
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create user",
			})
		}
		response := dto.UserResponse{
			Id: int(newUser.ID),
			Name: newUser.Name,
			Email: newUser.Email,
		}

		return  c.Status(fiber.StatusCreated).JSON(response)
	}
}
