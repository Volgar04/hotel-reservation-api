package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Volgar04/hotel-reservation/api"
	"github.com/Volgar04/hotel-reservation/db"
	"github.com/Volgar04/hotel-reservation/db/fixtures"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	var (
		ctx           = context.Background()
		mongoEndpoint = os.Getenv("MONGO_DB_URL")
		mongoDBName   = os.Getenv("MONGO_DB_NAME")
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(mongoDBName).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := db.Store{
		User:    db.NewMongoUserStore(client),
		Booking: db.NewMongoBookingStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   hotelStore,
	}

	user := fixtures.AddUser(&store, "james", "foo", false)
	fmt.Println("james ->", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(&store, "admin", "admin", true)
	fmt.Println("admin ->", api.CreateTokenFromUser(admin))
	hotel := fixtures.AddHotel(&store, "some hotel", "bermuda", 5, nil)
	room := fixtures.AddRoom(&store, "small", true, 99.9, hotel.ID)
	booking := fixtures.AddBooking(&store, user.ID, room.ID, 2, time.Now().AddDate(0, 1, 0), time.Now().AddDate(0, 1, 5), false)
	fmt.Println("booking ->", booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(&store, name, location, rand.Intn(5), nil)
	}
}
