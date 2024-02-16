package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/Volgar04/hotel-reservation/db"
	"github.com/Volgar04/hotel-reservation/types"
)

func AddBooking(store *db.Store, uid, rid primitive.ObjectID, numPersons int, fromDate, tillDate time.Time, canceled bool) *types.Booking {
	booking := types.Booking{
		UserID:     uid,
		RoomID:     rid,
		NumPersons: numPersons,
		FromDate:   fromDate,
		TillDate:   tillDate,
		Canceled:   canceled,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.TODO(), &booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, size string, ss bool, price float64, hid primitive.ObjectID) *types.Room {
	room := types.Room{
		Size:    size,
		Price:   price,
		HotelID: hid,
		Seaside: ss,
	}
	insertedRoom, err := store.Room.InsertRoom(context.TODO(), &room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHotel(store *db.Store, name string, loc string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if roomIDS == nil {
		roomIDS = []primitive.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		FirstName: fn,
		LastName:  ln,
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	insertedUser, err := store.User.InsertUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
