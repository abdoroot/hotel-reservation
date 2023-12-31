package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/abdoroot/hotel-reservation/api"
	"github.com/abdoroot/hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	DBURI := os.Getenv(db.MONGODBENVDBURI)
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(DBURI))
	if err != nil {
		panic(err)
	}

	var (
		app        = fiber.New(fiber.Config{ErrorHandler: api.ErrorHandler})
		us         = db.NewMongoUserStore(client)
		hs         = db.NewMongoHotelStore(client)
		rs         = db.NewMongoRoomStore(client, hs)
		bs         = db.NewMongoBookingStore(client)
		apiv1      = app.Group("/api/v1", api.JWTAuthentication(us))
		adminRoute = apiv1.Group("/admin", api.AdminAuth)
		authRouter = app.Group("/api/auth")
		store      = &db.Store{
			User:    us,
			Hotel:   hs,
			Room:    rs,
			Booking: bs,
		}
		hotelHandler   = api.NewHotelHandler(store)
		userHandler    = api.NewUserHandler(store)
		authHandler    = api.NewAuthHandler(store)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
	)

	// - Auth route
	authRouter.Post("/", authHandler.HandleAuthUser)
	// - user route
	apiv1.Delete("user/:id", userHandler.HandleDeleteUser) //update user
	apiv1.Put("user/:id", userHandler.HandlePutUser)       //update user
	apiv1.Post("user", userHandler.HandlePostUser)         //create user
	apiv1.Get("user/:id", userHandler.HandleGetUser)       //get userById
	apiv1.Get("users", userHandler.HandleGetUsers)         //get user
	// - hotel route
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGethotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// - room route
	apiv1.Post("/room/:id/book", roomHandler.HandleRoomBooking)

	// - bookings
	apiv1.Get("/booking/:id", bookingHandler.HandleGetbooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleGetCancelbooking)
	// - admin routes
	adminRoute.Get("/booking", bookingHandler.HandleGetbookings)

	listenAddr := os.Getenv("HTTP_LISTER_ADDRESS")
	if err := app.Listen(listenAddr); err != nil {
		fmt.Println("starting server err:", err)
	}
}
