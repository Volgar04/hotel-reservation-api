package api

import (
	"context"
	"log"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Volgar04/hotel-reservation/db"
	"github.com/joho/godotenv"
)

type testDB struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testDB) teardown(t *testing.T) {
	dbname := os.Getenv(db.MongoDBNameEnvName)
	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testDB {
	if err := godotenv.Load("../.env"); err != nil {
		t.Error(err)
	}
	dburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	return &testDB{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
			Hotel:   hotelStore,
		},
	}
}
