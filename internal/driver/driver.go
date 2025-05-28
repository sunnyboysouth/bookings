package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"        // Import the pq driver
	_ "github.com/jackc/pgx/v5"        // Import the pgx driver
	_ "github.com/jackc/pgx/v5/stdlib" // Import the pgx driver
)

// DB holds the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbConnLifetime = 5 * time.Minute // 0 means no limit

// ConnerectSQL connects to the database using the provided DSN (Data Source Name).
func ConnectSQL(dsn string) (*DB, error) {
	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbConnLifetime)

	dbConn.SQL = db

	err = testDB(db) // test the connection again! dont know why?!

	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

// tries to ping the database
func testDB(d *sql.DB) error {
	if err := d.Ping(); err != nil {
		return err
	}
	return nil
}

// NewDatabase creates a new database connection pool using the provided DSN (Data Source Name).
func NewDatabase(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
