// Generally following this: https://earthly.dev/blog/golang-sqlite/
package store

import (
	"database/sql"
	"os"
	"time"

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

func (d *DbStore) InsertRegistration(registrationData api.DonationRegistration) (*api.DonationRegistration, error) {
	// Write information to db
	now := time.Now()
	res, err := d.db.Exec("INSERT INTO donations VALUES (NULL,?,?,?,?,?)", now.String(), registrationData.Name, registrationData.Type, registrationData.Quantity, registrationData.Description)
	if err != nil {
		return nil, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return d.GetDonationRegistration(int(id))
}

func (d *DbStore) GetDonationRegistration(id int) (*api.DonationRegistration, error) {
	row := d.db.QueryRow("SELECT * FROM donations WHERE Id=?", id)
	registration := api.DonationRegistration{}
	if err := row.Scan(&registration.Id, &registration.Date, &registration.Name, &registration.Type, &registration.Quantity, &registration.Description); err == sql.ErrNoRows {
		return nil, err
	}

	return &registration, nil
}

func (d *DbStore) InsertDistribution(distributionData api.DonationDistribution) (*api.DonationDistribution, error) {
	now := time.Now()
	res, err := d.db.Exec("INSERT INTO donation_distributions VALUES (NULL,?,?,?,?,?)", distributionData.DonationId, now.String(), distributionData.Type, distributionData.Quantity, distributionData.Description)
	if err != nil {
		return nil, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return nil, err
	}
	return d.GetDonationDistribution(int(id))
}

func (d *DbStore) GetDonationDistribution(id int) (*api.DonationDistribution, error) {
	row := d.db.QueryRow("SELECT * FROM donation_distributions WHERE id=?", id)
	distribution := api.DonationDistribution{}
	if err := row.Scan(&distribution.Id, &distribution.DonationId, &distribution.Date, &distribution.Type, &distribution.Quantity, &distribution.Description); err == sql.ErrNoRows {
		return nil, err
	}

	return &distribution, nil
}

func (d *DbStore) GetDistributedDonationAmount(donationId int) (int64, error) {
	row := d.db.QueryRow("SELECT SUM(quantity) FROM donation_distributions WHERE donation_id=?", donationId)
	var amount int64
	if err := row.Scan(&amount); err == sql.ErrNoRows {
		return 0, err
	}
	return amount, nil
}

func (d *DbStore) GetInventoryReport() (*api.DonationInventory, error) {

	// For each type
	types := []api.DonationType{api.Clothing, api.Food, api.Money}
	var reportByType []api.TypeDonationStatus

	for _, t := range types {
		// For all donations for the type
		registrationRows, err := d.db.Query("SELECT * FROM donations WHERE type=?", t)
		var donationStatuses []api.DonationStatus
		if err != nil {
			return nil, err
		}
		defer registrationRows.Close()
		for registrationRows.Next() {
			registration := api.DonationRegistration{}
			err = registrationRows.Scan(&registration.Id, &registration.Date, &registration.Name, &registration.Type, &registration.Quantity, &registration.Description)
			if err != nil {
				return nil, err
			}

			// For all distributions for the donation
			distributionRows, err := d.db.Query("SELECT * FROM donation_distributions WHERE donation_id=?", registration.Id)
			if err != nil {
				return nil, err
			}
			defer distributionRows.Close()
			var distributions []api.DonationDistribution
			for distributionRows.Next() {
				distribution := api.DonationDistribution{}
				err = distributionRows.Scan(&distribution.Id, &distribution.DonationId, &distribution.Date, &distribution.Type, &distribution.Quantity, &distribution.Description)
				if err != nil {
					return nil, err
				}
				distributions = append(distributions, distribution)
			}
			status := api.DonationStatus{Donation: registration, Distributions: distributions}
			donationStatuses = append(donationStatuses, status)
		}
		reportByType = append(reportByType, api.TypeDonationStatus{Type: t, Statuses: donationStatuses})
	}
	report := api.DonationInventory{Report: &reportByType}

	return &report, nil
}

func (d *DbStore) GetDonorReport() (*api.DonorReport, error) {
	// Get unique donors
	nameRows, err := d.db.Query("SELECT DISTINCT(name) FROM donations")
	if err != nil {
		return nil, err
	}
	defer nameRows.Close()

	var donorSummaries []api.DonorSummary

	for nameRows.Next() {
		var name string
		err = nameRows.Scan(&name)
		if err != nil {
			return nil, err
		}

		var donationSummaries []api.DonationSummary

		// Aggregate donations by type
		// Aggregate distributions by type
		donationSummaryRows, err := d.db.Query(`
		WITH dr_summary AS (
			  SELECT type, SUM(quantity)
			    FROM donations
			   WHERE name = ?
			GROUP BY type
		),
		WITH dd_summary AS (
			SELECT type, SUM(quantity)
			  FROM donation_distributions
			 WHERE name = ?
		  GROUP BY type
		)
		SELECT dr_summary.type, dr_summary.quantity, dd_summary.quantity
		  FROM dr_summary
		  JOIN dd_summary
		    ON dr_summary.type = dd_summary.type;
		`, name, name)
		if err != nil {
			return nil, err
		}
		defer donationSummaryRows.Close()

		for donationSummaryRows.Next() {
			donation := api.DonationSummary{}
			donationSummaryRows.Scan(&donation.Type, &donation.Quantity, &donation.QuantityDistributed)
			donationSummaries = append(donationSummaries, donation)
		}

		donorSummary := api.DonorSummary{Donations: donationSummaries, Name: name}
		donorSummaries = append(donorSummaries, donorSummary)
	}

	report := api.DonorReport{Report: &donorSummaries}
	return &report, nil
}
