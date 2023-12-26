package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	jwtTokenExpireAfter = 4 //hours
)

func JWTAuthentication(store db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token, ok := c.GetReqHeaders()["X-Api-Key"]
		if !ok {
			return fmt.Errorf("not authorized")
		}
		claims, err := ParseToken(token[0])
		if err != nil {
			fmt.Println("ParseJWTToken :", err)
			return fmt.Errorf("not authorized")
		}

		userId := claims["user_id"].(string)
		oid, err := primitive.ObjectIDFromHex(userId)
		if err != nil {
			fmt.Println("fail to convert user id string to objectId :", err)
			return fmt.Errorf("not authorized")
		}

		user, err := store.GetUserByID(c.Context(), oid)
		if err != nil {
			fmt.Println("user not fount in db :", err)
			return fmt.Errorf("not authorized")
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
