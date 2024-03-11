package middeware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--- JWT authting")

	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	fmt.Println(token)
	if err := parseToken(token[0]); err != nil {
		return err
	}

	fmt.Println("token: ", token)

	return nil
}

func parseToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method ", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")

		fmt.Println("JWT_SECRET: ", secret)

		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("Failed to parse JWT token: ", err)
		return fmt.Errorf("unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("unauthorized")
}
