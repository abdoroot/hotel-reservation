package api

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const (
	jwtTokenExpireAfter = 4 //hours
)

func JWTAuthentication(store db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Key"]
		if !ok {
			return ErrorUnauthorized()
		}
		claims, err := ParseToken(token[0])
		if err != nil {
			fmt.Println("ParseJWTToken :", err)
			return ErrorUnauthorized()
		}

		userId := claims["user_id"].(string)
		user, err := store.GetUserByID(c.Context(), userId)
		if err != nil {
			fmt.Println("user not fount in db :", err)
			return ErrorUnauthorized()
		}

		//pass the auth user to context
		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func ParseToken(tokenString string) (jwt.MapClaims, error) {
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		fmt.Println("JWT_SECRET_KEY not set")
		return nil, NewError(http.StatusInternalServerError,"internal error")
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

	return nil, ErrorUnauthorized()
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
		return "", NewError(http.StatusInternalServerError,"Fail to create token")
	}
	return tokenString, nil
}
