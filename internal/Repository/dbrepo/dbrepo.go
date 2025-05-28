package dbrepo

import (
	"database/sql"

	repository "github.com/dsundar/bookings/internal/Repository"
	"github.com/dsundar/bookings/internal/config"
	"github.com/dsundar/bookings/internal/models"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// testDBRepo is a test implementation of the DatabaseRepo interface
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// UpdateUser is a mock implementation of the UpdateUser method
func (m *testDBRepo) UpdateUser(user models.User) error {
	// Mock implementation: pretend the update was successful
	return nil
}

// GetUserById is a mock implementation of the GetUserById method
func (m *testDBRepo) GetUserById(id int) (models.User, error) {
	// Mock implementation: return a dummy user
	return models.User{
		ID:        id,
		FirstName: "Test",
		LastName:  "User",
		Email:     "testuser@example.com",
	}, nil
}

func NewPostgresDBRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
