package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/Volgar04/hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup()
	defer db.teardown(t)

	user := fixtures.AddUser(db.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(db.Store, "hotel1", "a", 4, nil)
	room := fixtures.AddRoom(db.Store, "small", true, 100, hotel.ID)

	from := time.Now().AddDate(0, 1, 0)
	till := time.Now().AddDate(0, 1, 5)
	booking := fixtures.AddBooking(db.Store, user.ID, room.ID, 2, from, till, false)
	fmt.Println(booking)
}
