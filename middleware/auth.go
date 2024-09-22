package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
)

func Auth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	if authHeaderParts[0] != "Bearer" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	accessToken := authHeaderParts[1]
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		}

		jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
		return []byte(jwtSecretKey), nil
	})

	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthenticated",
		})
	}

	c.Locals("userInfo", claims)
	return c.Next()
}
