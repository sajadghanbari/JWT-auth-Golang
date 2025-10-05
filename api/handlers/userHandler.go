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
		var req dto.CreateUserRequest
		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse json",
			})
		}

		err = validate.Struct(req)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		hashedPassword := services.HashPassword(req.Password)
		newUser := models.User{
			Name:     req.Name,
			Email:    req.Email,
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
			Id:    int(newUser.ID),
			Name:  newUser.Name,
			Email: newUser.Email,
		}

		return c.Status(fiber.StatusCreated).JSON(response)
	}
}

func GetAllUsers(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var users []models.User
		//get all users from DB
		err := db.Find(&users).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not fetch users",
			})
		}

		var userResponses []dto.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, dto.UserResponse{
				Id:    int(user.ID),
				Name:  user.Name,
				Email: user.Email,
			})

		}

		return c.Status(fiber.StatusOK).JSON(userResponses)
	}
}

func DeleteUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var user models.User

		err := db.First(&user, id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not find user",
			})
		}

		err = db.Delete(&user).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not delete user",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "User deleted successfully",
		})
	}
}

func UpdateUser(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")

		var user models.User
		err := db.First(&user, id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"error": "User not found",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not find user",
			})
		}

		var req dto.UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}


		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		//TODO this part can be better
		if req.Name != nil {
			user.Name = *req.Name
		}
		if req.Email != nil {
			user.Email = *req.Email
		}
		if req.Password != nil {
			user.Password = services.HashPassword(*req.Password)
		}

		if err := db.Save(&user).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not update user",
			})
		}

		response := dto.UserResponse{
			Id:    int(user.ID),
			Name:  user.Name,
			Email: user.Email,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}
}
