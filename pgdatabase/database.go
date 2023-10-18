package pgdatabase

import (
	"database/sql"
	"fmt"
	"log"

	// Import the pq driver to connect to PostgreSQL
	_ "github.com/lib/pq"
)

// DB is the database connection object
var DB *sql.DB

// Config holds the database connection parameters
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// InitDB initializes the database with the given configuration
func InitDB(cfg Config) error {
	// Create the connection string using the given configuration
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName,
	)

	// Attempt to connect to the database
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Ping the database to verify the connection
	if err := DB.Ping(); err != nil {
		return err
	}

	log.Println("Connected to the database successfully.")
	return nil
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
