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
	for i := range 5 {
		db, err = sql.Open("postgres", connString)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to PostgreSQL!")
				break
			}
		}
		log.Printf("Waiting for database to be ready... (%d/5) %v\n", i+1, err)
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

func (r *PostgreSQLRepository) InsertVehicleLocation(vehicleLocation models.VehicleLocation) error {
	query := `INSERT INTO vehicle_locations (vehicle_id, latitude, longitude, timestamp) VALUES ($1, $2, $3, $4)`
	_, err := r.DB.Exec(query, vehicleLocation.VehicleID, vehicleLocation.Latitude, vehicleLocation.Longitude, vehicleLocation.Timestamp)
	if err != nil {
		return fmt.Errorf("failed to insert vehicle location: %w", err)
	}
	return nil
}

func (r *PostgreSQLRepository) GetLastVehicleLocation(vehicleID string) (models.VehicleLocation, error) {
	var vehicleLocation models.VehicleLocation
	query := `SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 ORDER BY timestamp DESC LIMIT 1`
	row := r.DB.QueryRow(query, vehicleID)
	err := row.Scan(&vehicleLocation.VehicleID, &vehicleLocation.Latitude, &vehicleLocation.Longitude, &vehicleLocation.Timestamp)
	if err == sql.ErrNoRows {
		return vehicleLocation, fmt.Errorf("no location found for vehicle %s", vehicleID)
	} else if err != nil {
		return vehicleLocation, fmt.Errorf("failed to get last location: %w", err)
	}
	return vehicleLocation, nil
}

func (r *PostgreSQLRepository) GetVehicleLocationHistory(vehicleID string, start, end int64) ([]models.VehicleLocation, error) {
	var vehicleLocations []models.VehicleLocation
	query := `SELECT vehicle_id, latitude, longitude, timestamp FROM vehicle_locations WHERE vehicle_id = $1 AND timestamp >= $2 AND timestamp <= $3 ORDER BY timestamp ASC`
	rows, err := r.DB.Query(query, vehicleID, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to get vehicle location history %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var vehicleLocation models.VehicleLocation
		if err := rows.Scan(&vehicleLocation.VehicleID, &vehicleLocation.Latitude, &vehicleLocation.Longitude, &vehicleLocation.Timestamp); err != nil {
			return nil, fmt.Errorf("failed to scan vehicle location hystory row: %w", err)
		}
		vehicleLocations = append(vehicleLocations, vehicleLocation)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during history iteration: %w", err)
	}

	return vehicleLocations, nil
}
