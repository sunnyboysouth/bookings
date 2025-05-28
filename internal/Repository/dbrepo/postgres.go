package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/dsundar/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *PostgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the database
func (m *PostgresDBRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	var newID int
	// insert the reservation into the database
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, user_id, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id`

	// replaced ExecContext with QueryRowContext which only returns err
	//QueryRowContext is used to execute a query that returns a single row and scan the result into a variable
	err := m.DB.QueryRowContext(ctx, stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.UserId,
		time.Now(),
		time.Now(),
	).Scan(&newID) // scan the returned id into newID

	if err != nil {
		log.Println("Error inserting reservation: ", err)
		return 0, err
	}

	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (m *PostgresDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	stmt := `insert into room_restriction (start_date, end_date, room_id, reservation_id, restriction_id, user_id, created_at, updated_at)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		r.UserId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		log.Println("Error inserting room restriction: ", err)
		return err
	}

	return nil
}

// SearchAvailabilityByDatesByRoomId checks if a room is available for the given dates
// It returns true if the room is available, false otherwise
func (m *PostgresDBRepo) SearchAvailabilityByDatesByRoomId(start, end time.Time, roomID int) (bool, error) {

	var numRows int

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	query := `select count(id) from room_restriction where room_id = $1 and $2 < end_date and $3 > start_date`

	row := m.DB.QueryRowContext(ctx, query, roomID, start, end)
	err := row.Scan(&numRows)

	if err != nil {
		log.Println("Error scanning row: ", err)
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// SearchAvailabilityForAllRooms checks if any room is available for the given date range and returns a slice of available rooms
func (m *PostgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	var rooms []models.Room

	query := `select 
		r.id, r.room_name
	from 
		rooms r
	where
		r.id not in 
	(select rr.room_id from room_restriction rr where $1 < rr.end_date and $2 > rr.start_date)`

	rows, err := m.DB.QueryContext(ctx, query, start, end)
	if err != nil {
		log.Println("Error querying rooms: ", err)
		return rooms, err
	}

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			log.Println("Error scanning room: ", err)
			return rooms, err
		}

		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil {
		log.Println("Error iterating over rows: ", err)
		return rooms, err
	}
	defer rows.Close()

	return rooms, nil
}

// GetRoomById retrieves a room by its ID
func (m *PostgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(&room.ID, &room.RoomName, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		log.Println("Error scanning room: ", err)
		return room, err
	}

	return room, nil
}

// GetUserById retrieves a user by their ID
func (m *PostgresDBRepo) GetUserById(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	var user models.User

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.AccessLevel, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Error scanning user: ", err)
		return user, err
	}

	return user, nil
}

// UpdateUser updates a user's information in the database
func (m *PostgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, email = $3, password = $4, access_level = $5, updated_at = $6 where id = $7`

	_, err := m.DB.ExecContext(ctx, query,
		u.FirstName,
		u.LastName,
		u.Email,
		u.Password,
		u.AccessLevel,
		time.Now(),
		u.ID,
	)
	if err != nil {
		log.Println("Error updating user: ", err)
		return err
	}

	return nil
}

// Authenticate checks if the user exists and if the password matches
func (m *PostgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // create a context with a timeout
	defer cancel()

	var id int
	var hashedPassword string

	query := `select id, password from users where email = $1`
	row := m.DB.QueryRowContext(ctx, query, email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.Println("Password mismatch for user:", email)
			return 0, "", errors.New("incorrect Passworrd") // return 0 and empty string if password does not match
		}
		log.Println("Error comparing password hash:", err)
		return 0, "", err // return error if there is an issue with comparing the hash
	}

	return id, hashedPassword, nil // return user ID and hashed password if authentication is successful
}
