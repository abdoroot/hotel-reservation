package types

import (
	"fmt"
	"log"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 4
)

type UpdateRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type User struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName         string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Email             string             `json:"email" bson:"email"`
	EncreptedPassword string             `json:"-" bson:"encrepted_password"`
}

type CreateUserRequest struct {
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	Email             string `json:"email"`
	EncreptedPassword string `json:"password"`
}

func (u UpdateRequest) ToBSON() bson.M {
	return bson.M{
		"first_name": u.FirstName,
		"last_name":  u.LastName,
	}
}

func (c CreateUserRequest) CreateUserFromUserRequest() (*User, error) {
	EncreptedPassword, err := bcrypt.GenerateFromPassword([]byte(c.EncreptedPassword), bcryptCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		FirstName:         c.FirstName,
		LastName:          c.LastName,
		Email:             c.Email,
		EncreptedPassword: string(EncreptedPassword),
	}
	return user, nil
}

func (c CreateUserRequest) Validate() []error {
	err := make([]error, 0)
	if len(c.FirstName) < 2 {
		err = append(err, fmt.Errorf("First name must be at least 2 char"))
	}
	if len(c.LastName) < 2 {
		err = append(err, fmt.Errorf("Last name must be at least 2 char"))
	}
	if !strings.Contains(c.Email, "@") {
		err = append(err, fmt.Errorf("Email not valid"))
	}
	if len(c.EncreptedPassword) < 4 {
		log.Println(c.EncreptedPassword)
		err = append(err, fmt.Errorf("Password length must be at least 4 char"))
	}
	return err
}

func (u UpdateRequest) Validate() []error {
	err := make([]error, 0)
	if len(u.FirstName) < 2 {
		err = append(err, fmt.Errorf("First name must be at least 2 char"))
	}
	if len(u.LastName) < 2 {
		err = append(err, fmt.Errorf("Last name must be at least 2 char"))
	}
	return err
}
