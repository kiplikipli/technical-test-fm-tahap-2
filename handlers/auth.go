package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kiplikipli/technical-test-fm-tahap-2/services"
	"gorm.io/gorm"
)

type (
	RegisterRequest struct {
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
		PhoneNumber string `json:"phone_number" validate:"required"`
		Address     string `json:"address" validate:"required"`
		Pin         string `json:"pin" validate:"required"`
	}

	RegisterResponse struct {
		UserID      string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		CreatedAt   string `json:"created_at"`
	}

	LoginRequest struct {
		PhoneNumber string `json:"phone_number" validate:"required"`
		Pin         string `json:"pin" validate:"required"`
	}

	LoginResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	UpdateProfileRequest struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Address   string `json:"address"`
	}

	UpdateProfileResponse struct {
		UserID      string `json:"user_id"`
		FirstName   string `json:"first_name"`
		LastName    string `json:"last_name"`
		PhoneNumber string `json:"phone_number"`
		Address     string `json:"address"`
		UpdatedDate string `json:"updated_date"`
	}
)

func Register(c *fiber.Ctx) error {
	json := new(RegisterRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	_, err := services.GetUserByPhoneNumber(json.PhoneNumber)
	if err != gorm.ErrRecordNotFound {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Phone Number already registered",
		})
	}

	newUser, err := services.CreateUser(&services.User{
		FirstName:   json.FirstName,
		LastName:    json.LastName,
		PhoneNumber: json.PhoneNumber,
		Address:     json.Address,
		Pin:         json.Pin,
	})

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "SUCCESS",
		"result": RegisterResponse{
			UserID:      newUser.ID.String(),
			FirstName:   newUser.FirstName,
			LastName:    newUser.LastName,
			PhoneNumber: newUser.PhoneNumber,
			Address:     newUser.Address,
			CreatedAt:   newUser.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func Login(c *fiber.Ctx) error {
	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	user, err := services.GetUserByPhoneNumber(json.PhoneNumber)
	if err == gorm.ErrRecordNotFound {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Phone Number and PIN doesn't match",
		})
	}

	err = services.CompareHash(user.Pin, json.Pin)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Phone Number and PIN doesn't match",
		})
	}

	jwtClaims := jwt.MapClaims{
		"iss":     "technical-test-fm-tahap-2",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
		"user_id": user.ID.String(),
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	accessToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status": "SUCCESS",
		"result": LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: "refresh_token",
		},
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	json := new(UpdateProfileRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	jwtClaims := c.Locals("userInfo").(jwt.MapClaims)
	userId := jwtClaims["user_id"].(string)
	userUuid, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Invalid UUID",
		})
	}

	updateRequest := services.User{
		ID:        userUuid,
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Address:   json.Address,
	}

	updatedUser, err := services.UpdateUser(&updateRequest)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "SUCCESS",
		"result": updatedUser,
	})
}
