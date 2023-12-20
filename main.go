package main

import (
	"context"

	"github.com/abdoroot/hotel-reservation/api"
	"github.com/abdoroot/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongoUri    = "mongodb://localhost:27017/?tls=false"
	mongodbName = "hotel-reservation"
)

func main() {
	app := fiber.New()
	apiv1 := app.Group("/api")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	userStore := db.NewMongoStore(client, mongodbName)
	userHandler := api.NewUserHandler(userStore)
	{
		apiv1.Delete("user/:id", userHandler.HandleDeleteUser) //update user
		apiv1.Put("user/:id", userHandler.HandlePutUser)       //update user
		apiv1.Post("user", userHandler.HandlePostUser)         //create user
		apiv1.Get("user/:id", userHandler.HandleGetUser)       //get userById
		apiv1.Get("users", userHandler.HandleGetUsers)         //get user
	}

	app.Listen(":3000")
}
