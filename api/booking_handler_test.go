package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Volgar04/hotel-reservation/db/fixtures"
	"github.com/Volgar04/hotel-reservation/types"
)

func TestUserGetBooking(t *testing.T) {
	db := setup()
	defer db.teardown(t)

	var (
		otherUser      = fixtures.AddUser(db.Store, "jimmy", "foo", false)
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		hotel          = fixtures.AddHotel(db.Store, "hotel1", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 100, hotel.ID)
		from           = time.Now().AddDate(0, 1, 0)
		till           = time.Now().AddDate(0, 1, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, 2, from, till, false)
		app            = fiber.New()
		route          = app.Group("/", JWTAuthentication(db.User))
		bookingHandler = NewBookingHandler(db.Store)
	)
	route.Get("/:id", bookingHandler.HandleGetBooking)
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	var bookings *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if bookings.ID != booking.ID {
		t.Errorf("expected booking id %s, got %s", booking.ID, bookings.ID)
	}
	if bookings.UserID != user.ID {
		t.Errorf("expected user id %s, got %s", user.ID, bookings.UserID)
	}

	// test non-owner cannot access the booking
	req = httptest.NewRequest("GET", fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(otherUser))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected status code 500, got %d", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	db := setup()
	defer db.teardown(t)

	var (
		user           = fixtures.AddUser(db.Store, "james", "foo", false)
		adminUser      = fixtures.AddUser(db.Store, "admin", "admin", true)
		hotel          = fixtures.AddHotel(db.Store, "hotel1", "a", 4, nil)
		room           = fixtures.AddRoom(db.Store, "small", true, 100, hotel.ID)
		from           = time.Now().AddDate(0, 1, 0)
		till           = time.Now().AddDate(0, 1, 5)
		booking        = fixtures.AddBooking(db.Store, user.ID, room.ID, 2, from, till, false)
		app            = fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
		admin          = app.Group("/", JWTAuthentication(db.User), AdminAuth)
		bookingHandler = NewBookingHandler(db.Store)
	)
	admin.Get("/", bookingHandler.HandleGetBookings)
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	have := bookings[0]
	if have.ID != booking.ID {
		t.Errorf("expected booking id %d, got %d", booking.ID, have.ID)
	}
	if have.UserID != user.ID {
		t.Errorf("expected user id %d, got %d", user.ID, have.UserID)
	}

	// test non-admin cannot access the bookings
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))
	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status unauthorized, got %d", resp.StatusCode)
	}
}
