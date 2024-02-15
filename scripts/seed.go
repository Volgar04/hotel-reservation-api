package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Volgar04/hotel-reservation/db"
	"github.com/Volgar04/hotel-reservation/types"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(isAdmin bool, fname, lname, email, password string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lname,
		Email:     email,
		Password:  password,
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = isAdmin
	_, err = userStore.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	rooms := []types.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "suite",
			Price: 199.9,
		},
		{
			Size:  "normal",
			Price: 159.9,
		},
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		_, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Beluccia", "France", 4)
	seedHotel("The Cozy Hotel", "Nederlands", 3)
	seedHotel("Dont die in your sleep", "London", 1)
	seedUser(false, "james", "foo", "james@foo.com", "supersecurepassword")
	seedUser(true, "nicolas", "martin", "nicolas@martin.com", "adminsupersecurepassword")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
