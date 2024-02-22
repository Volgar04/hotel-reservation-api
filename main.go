package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Volgar04/hotel-reservation/api"
	"github.com/Volgar04/hotel-reservation/db"
)

// Configuration
// 1. MongoDB endpoint
// 2. ListenAddress of our HTTP server
// 3. JWT secret
// 4. MongoDBName

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}

	// handlers initializations
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
		}
		userHandler    = api.NewUserHandler(userStore)
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)
		app            = fiber.New(config)
		auth           = app.Group("/api/")
		apiV1          = app.Group("/api/v1", api.JWTAuthentication(userStore))
		admin          = apiV1.Group("/admin", api.AdminAuth)
	)

	// auth
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// Versioned API routes
	// user handlers
	apiV1.Post("/user", userHandler.HandlePostUser)
	apiV1.Get("/user", userHandler.HandleGetUsers)
	apiV1.Get("/user/:id", userHandler.HandleGetUser)
	apiV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiV1.Put("/user/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiV1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiV1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// rooms handlers
	apiV1.Get("room", roomHandler.HandleGetRooms)
	apiV1.Post("room/:id/book", roomHandler.HandleBookRoom)
	// TODO: cancel a booking

	// bookings handlers
	apiV1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiV1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin handlers
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	err = app.Listen(listenAddr)
	if err != nil {
		panic(err)
	}
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
