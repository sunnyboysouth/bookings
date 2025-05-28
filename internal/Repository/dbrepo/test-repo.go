package dbrepo

import (
	"errors"
	"time"

	"github.com/dsundar/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	// if the room id is 2, return an error
	if res.RoomID == 2 {
		return 0, errors.New("room id 2 is not available")
	}
	return 1, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 1000 {
		return errors.New("room id 1000 is not available")
	}
	return nil
}

// SearchAvailabilityByDatesByRoomId checks if a room is available for the given dates
// It returns true if the room is available, false otherwise
func (m *testDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error) {
	if roomID == 2 {
		return true, errors.New("room id 2 is not available")
	}

	return false, nil
}

// SearchAvailabilityForAllRooms checks if any room is available for the given date range and returns a slice of available rooms
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	if start.Year() == 2020 {
		return rooms, errors.New("rooms not available")
	}
	return rooms, nil
}

// GetRoomById retrieves a room by its ID
func (m *testDBRepo) GetRoomById(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("room not found")
	}

	return room, nil
}

func (m *testDBRepo) Authenticate(email, password string) (int, string, error) {
	return 1, "user", nil
}
