package repository

import (
	"database/sql"
	"fmt"
	"log"
	"sistem-manajemen-armada/internal/models"
	"time"

	_ "github.com/lib/pq"
)

type PostgreSQLRepository struct {
	DB *sql.DB
}

func NewPostgreSQLRepository(host, port, user, password, dbname string) (*PostgreSQLRepository, error) {
	connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	// Retry connection
	for i := range 10 {
		db, err = sql.Open("postgres", connString)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to PostgreSQL!")
				break
			}
		}
		log.Printf("Waiting for database to be ready... (%d/10) %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		return nil, fmt.Errorf("could not connect to PostgreSQL after multiple attempts: %w", err)
	}

	repo := &PostgreSQLRepository{DB: db}
	if err := repo.createTable(); err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}
	return repo, nil
}

func (r *PostgreSQLRepository) createTable() error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS vehicle_locations (
		id SERIAL PRIMARY KEY,
		vehicle_id VARCHAR(255) NOT NULL,
		latitude DOUBLE PRECISION NOT NULL,
		longitude DOUBLE PRECISION NOT NULL,
		timestamp BIGINT NOT NULL
	);`

	_, err := r.DB.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("error creating table: %w", err)
	}
	log.Println("Table 'vehicle_locations' checked/created successfully.")
	return nil
}

func (r *PostgreSQLRepository) InsertLocation(loc models.VehicleLocation) error {
	query := `INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, loc.VehicleID, loc.Latitude, loc.Longitude, loc.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert location: %w", err)
	}
	return nil
}
