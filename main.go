package main

import (
	"context"

	"github.com/abdoroot/hotel-reservation/api"
	"github.com/abdoroot/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}

	var (
		app          = fiber.New()
		apiv1        = app.Group("/api")
		userStore    = db.NewMongoUserStore(client, db.DBName)
		userHandler  = api.NewUserHandler(userStore)
		hs           = db.NewMongoHotelStore(client)
		rs           = db.NewMongoRoomStore(client, hs)
		hotelHandler = api.NewHotelHandler(hs, rs)
	)
	{
		//user route
		apiv1.Delete("user/:id", userHandler.HandleDeleteUser) //update user
		apiv1.Put("user/:id", userHandler.HandlePutUser)       //update user
		apiv1.Post("user", userHandler.HandlePostUser)         //create user
		apiv1.Get("user/:id", userHandler.HandleGetUser)       //get userById
		apiv1.Get("users", userHandler.HandleGetUser)          //get user
	}

	{
		//hotel route
		apiv1.Get("/hotel", hotelHandler.HandleGetHotel)
		apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	}

	app.Listen(":3000")
}
