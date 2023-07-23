package store

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/nowei/open-seattle-example/server/internal/api"
	"github.com/nowei/open-seattle-example/server/internal/logger"
)

var log = logger.GetLogger().Sugar()

type DbStore struct {
	db *sql.DB
}

const file string = "shelter.db"
const schemaFile string = "schemas.sql"

func InstantiateDbStore() *DbStore {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("There was a problem opening the db connection: %s", err.Error())
		return nil
	}

	schemaBytes, err := os.ReadFile(schemaFile)
	if err != nil {
		log.Fatalf("There was a problem with reading the file: %s", err.Error())
		return nil
	}
	schema := string(schemaBytes)

	if _, err := db.Exec(schema); err != nil {
		log.Fatalf("Creating the schema: %s", err.Error())
		return nil
	}

	return &DbStore{
		db: db,
	}
}

func (d *DbStore) InsertRegistration(registration api.DonationRegistration) error {
	// Write information to db
	return nil
}

func (d *DbStore) InsertDistribution(distribution api.DonationDistribution) error {
	// Get donation

	// Aggregate all distributions of this donation

	// If new distribution + aggregate < original donation, add a new distribution

	return nil
}

func (d *DbStore) GetInventoryReport() (*api.DonationInventory, error) {

	// For each type

	// For all donations for the type
	// For all distributions for the donation
	return nil, nil
}

func (d *DbStore) GetDonorReport() (*api.DonorReport, error) {
	// For each donor

	// Aggregate donations by type

	// Aggregate distributions by type

	return nil, nil
}
