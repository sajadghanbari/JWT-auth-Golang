package controllers

import (
	"JWT-Authentication-go/database"
	"JWT-Authentication-go/models"
	"fmt"
	// "go/token"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

func Login(c *fiber.Ctx) error{
	fmt.Println("Received a login request")
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}


	var user models.User
	database.DB.Where("email = ?", data["email"]).First(&user)
	if user.ID == 0{
		fmt.Println("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data["password"]))
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   strconv.Itoa(int(user.ID)),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := claims.SignedString([]byte("secretKey"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	cookie := fiber.Cookie{
		Name:  "jwt",
		Value: token,
		Expires: time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
		Secure: true,
	}

	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful"})
}


func User(c *fiber.Ctx) error{
	fmt.Println("Request to get user")
	cookie := c.Cookies("jwt")

    token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte("secretKey"), nil
    })

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error" : "Unauthorized",
		})
	}

	claims , ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error" : "Failed to parse claim",
		})
	}

	id , _ := strconv.Atoi((*claims)["sub"].(string))
	user := models.User{ID : uint(id)}

	database.DB.Where("id=?",id).First(&user)

	return  c.JSON(user)
}

func Logout(c *fiber.Ctx) error{
	fmt.Println("Received a logout request")
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout successful"})
}