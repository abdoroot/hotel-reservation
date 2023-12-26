package main

import (
	"context"
	"fmt"

	"github.com/abdoroot/hotel-reservation/api"
	"github.com/abdoroot/hotel-reservation/db"
	"github.com/abdoroot/hotel-reservation/middleware"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(ctx *fiber.Ctx, err error) error {
		return ctx.JSON(bson.M{
			"error": err.Error(),
		})
	},
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		panic(err)
	}

	var (
		app        = fiber.New(config)
		us         = db.NewMongoUserStore(client)
		hs         = db.NewMongoHotelStore(client)
		rs         = db.NewMongoRoomStore(client, hs)
		bs         = db.NewMongoBookingStore(client)
		apiv1      = app.Group("/api/v1", middleware.JWTAuthentication(us))
		authRouter = app.Group("/api/auth")
		store      = &db.Store{
			User:    us,
			Hotel:   hs,
			Room:    rs,
			Booking: bs,
		}
		hotelHandler = api.NewHotelHandler(store)
		userHandler  = api.NewUserHandler(store)
		authHandler  = api.NewAuthHandler(store)
		roomHandler  = api.NewRoomHandler(store)
	)

	//Auth route
	authRouter.Post("/", authHandler.HandleAuthUser)
	//user route
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser) //update user
	apiv1.Put("user/:id", userHandler.HandlePutUser)       //update user
	apiv1.Post("user", userHandler.HandlePostUser)         //create user
	apiv1.Get("user/:id", userHandler.HandleGetUser)       //get userById
	apiv1.Get("users", userHandler.HandleGetUsers)         //get user
	//hotel route
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGethotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	//room route
	apiv1.Post("/room/:id/book", roomHandler.HandleRoomBooking)

	if err := app.Listen(":3000"); err != nil {
		fmt.Println("starting err:", err)
	}
}
