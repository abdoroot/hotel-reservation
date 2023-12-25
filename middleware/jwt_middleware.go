package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtTokenExpireAfter = 4 //hours
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Key"]
	if !ok {
		return fmt.Errorf("not authorized")
	}
	claims, err := ParseToken(token[0])
	fmt.Println(claims)
	if err != nil {
		fmt.Println("ParseJWTToken :", err)
		return fmt.Errorf("not authorized")
	}
	return c.Next()
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		fmt.Println("JWT_SECRET_KEY not set")
		return nil, fmt.Errorf("internal error")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func CreateUserJwtToken(user *types.User) (string, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	tokenExpire := time.Now().Add(jwtTokenExpireAfter * time.Hour).Unix()
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     tokenExpire,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Fail to create token")
	}
	return tokenString, nil
}
