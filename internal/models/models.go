package models

import "time"

// User is the struct for the users model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   string
	UpdatedAt   string
}

// Room is the struct for the rooms model
type Room struct {
	ID        int
	RoomName  string
	CreatedAt string
	UpdatedAt string
	UserId    int
}

// Restriction sis the struct for the restrictions model
type Restriction struct {
	ID               int
	RestrictionsName string
	CreatedAt        string
	UpdatedAt        string
	UserId           int
}

// Resevations holds reservations model
type Reservation struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomID    int
	UserId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

// RoomRestriction holds room restrictions model
type RoomRestriction struct {
	ID            int
	StartDate     time.Time
	EndDate       time.Time
	RoomID        int
	RestrictionID int
	ReservationID int
	UserId        int
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Reservation   Reservation
	Restriction   Restriction
}
