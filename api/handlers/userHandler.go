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


// CreateUser godoc
// @Summary      Create a new user
// @Description  Create a new user with name, email, and password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body  dto.CreateUserRequest  true  "User data"
// @Success      201  {object}  dto.UserResponse
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Router       /users/create [post]
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


// GetAllUsers godoc
// @Summary      Get all users
// @Description  Retrieve all users from the database
// @Tags         users
// @Produce      json
// @Success      200  {array}  dto.UserResponse
// @Failure      500  {object}  map[string]string
// @Router       /users/get-users [get]
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


// DeleteUser godoc
// @Summary      Delete a user
// @Description  Delete a user by their ID
// @Tags         users
// @Produce      json
// @Param        id   path   int  true  "User ID"
// @Success      200  {object}  map[string]string  "User deleted successfully"
// @Failure      404  {object}  map[string]string  "User not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /users/delete/{id} [delete]
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


// UpdateUser godoc
// @Summary      Update an existing user
// @Description  Update user details (only provided fields will be changed)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id    path   int                      true  "User ID"
// @Param        user  body   dto.UpdateUserRequest     true  "Updated user data"
// @Success      200   {object}  dto.UserResponse
// @Failure      400   {object}  map[string]string  "Invalid request or validation error"
// @Failure      404   {object}  map[string]string  "User not found"
// @Failure      500   {object}  map[string]string  "Internal server error"
// @Router       /users/update/{id} [put]
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
