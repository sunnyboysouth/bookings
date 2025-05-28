package repository

import (
	"time"

	"github.com/dsundar/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)

	InsertRoomRestriction(r models.RoomRestriction) error

	SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error)

	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)

	GetRoomById(id int) (models.Room, error)

	GetUserById(id int) (models.User, error)

	UpdateUser(u models.User) error

	Authenticate(email, password string) (int, string, error)
}
