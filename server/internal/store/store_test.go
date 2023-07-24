package store

import (
	"database/sql"
	"testing"

	"github.com/nowei/open-seattle-example/server/internal/api"
	"github.com/nowei/open-seattle-example/server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

const DBConnection = "file:test.db?mode=memory&cache=shared"
const TestSchemaFile = "schemas.sql"

func TestingStoreSetup() *DbStore {
	db, err := sql.Open("sqlite3", DBConnection)
	if err != nil {
		log.Fatalf("could not open db: %s", err.Error())
	}

	dbStore := InstantiateDbStore(db, TestSchemaFile)
	if dbStore == nil {
		log.Fatalf("could not create db")
	}
	return dbStore
}

func TestInsertDonation(t *testing.T) {
	t.Run("test add multiple registrations", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		registrationData := testutils.CreateTestDonationRegistration("Andrew", 42, api.Clothing, &description)
		registration, err := dbstore.InsertRegistration(registrationData)
		assert.NoError(t, err)
		assert.Equal(t, *registration.Description, description)
		assert.Equal(t, *registration.Id, 1)
		assert.Equal(t, registration.Name, "Andrew")
		assert.Equal(t, registration.Type, api.Clothing)
		assert.Equal(t, registration.Quantity, 42)

		// Test getting the donation registration
		registration, err = dbstore.GetDonationRegistration(1)
		assert.NoError(t, err)
		assert.Equal(t, *registration.Description, description)
		assert.Equal(t, *registration.Id, 1)
		assert.Equal(t, registration.Name, "Andrew")
		assert.Equal(t, registration.Type, api.Clothing)
		assert.Equal(t, registration.Quantity, 42)

		registrationData = testutils.CreateTestDonationRegistration("Andrew", 43, api.Money, &description)
		registration, err = dbstore.InsertRegistration(registrationData)
		assert.NoError(t, err)
		assert.Equal(t, *registration.Description, description)
		assert.Equal(t, *registration.Id, 2)
		assert.Equal(t, registration.Name, "Andrew")
		assert.Equal(t, registration.Type, api.Money)
		assert.Equal(t, registration.Quantity, 43)
	})
}

func TestInsertDistribtution(t *testing.T) {
	t.Run("test add multiple distributions", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		registrationData := testutils.CreateTestDonationRegistration("Andrew", 42, api.Clothing, &description)
		_, err := dbstore.InsertRegistration(registrationData)
		assert.NoError(t, err)

		distributionData := testutils.CreateTestDonationDistribution(1, 4, api.Clothing, &description)
		distribution, err := dbstore.InsertDistribution(distributionData)
		assert.NoError(t, err)
		assert.Equal(t, *distribution.Description, description)
		assert.Equal(t, *distribution.Id, 1)
		assert.Equal(t, distribution.DonationId, 1)
		assert.Equal(t, distribution.Type, api.Clothing)
		assert.Equal(t, distribution.Quantity, 4)

		distribution, err = dbstore.InsertDistribution(distributionData)
		assert.NoError(t, err)
		assert.Equal(t, *distribution.Description, description)
		assert.Equal(t, *distribution.Id, 2)
		assert.Equal(t, distribution.DonationId, 1)
		assert.Equal(t, distribution.Type, api.Clothing)
		assert.Equal(t, distribution.Quantity, 4)

		// Test getting the donation distribution
		distribution, err = dbstore.GetDonationDistribution(2)
		assert.NoError(t, err)
		assert.Equal(t, *distribution.Description, description)
		assert.Equal(t, *distribution.Id, 2)
		assert.Equal(t, distribution.DonationId, 1)
		assert.Equal(t, distribution.Type, api.Clothing)
		assert.Equal(t, distribution.Quantity, 4)

		// Test getting the distributed donation amount
		amount, err := dbstore.GetDistributedDonationAmount(1)
		assert.NoError(t, err)
		assert.Equal(t, amount, 8)
	})

	t.Run("test create donation distribution fails", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		distributionData := testutils.CreateTestDonationDistribution(5, 4, api.Clothing, &description)
		d, err := dbstore.InsertDistribution(distributionData)
		log.Infof("This is %v", d)
		assert.Error(t, err)
	})

}

func TestInventoryReport(t *testing.T) {
	t.Run("test single item type", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		registrationData := testutils.CreateTestDonationRegistration("Andrew", 42, api.Clothing, &description)
		_, _ = dbstore.InsertRegistration(registrationData)

		distributionData := testutils.CreateTestDonationDistribution(1, 4, api.Clothing, &description)
		_, _ = dbstore.InsertDistribution(distributionData)
		_, _ = dbstore.InsertDistribution(distributionData)

		report, err := dbstore.GetInventoryReport()
		assert.NoError(t, err)
		assert.NotNil(t, report.Report)
		for _, r := range *report.Report {
			if r.Type == api.Clothing {
				assert.NotNil(t, r.Statuses)
				assert.Equal(t, len(r.Statuses), 1)
				assert.Equal(t, len(r.Statuses[0].Distributions), 2)
				assert.Equal(t, *r.Statuses[0].Donation.Id, 1)
			} else {
				assert.Nil(t, r.Statuses)
			}
		}
	})

	t.Run("test all item types", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		for i, v := range []api.DonationType{api.Clothing, api.Food, api.Money} {
			registrationData := testutils.CreateTestDonationRegistration("Andrew", 42, v, &description)
			_, _ = dbstore.InsertRegistration(registrationData)
			distributionData := testutils.CreateTestDonationDistribution(i+1, 4, v, &description)
			_, _ = dbstore.InsertDistribution(distributionData)
			_, _ = dbstore.InsertDistribution(distributionData)
		}

		report, err := dbstore.GetInventoryReport()
		assert.NoError(t, err)
		assert.NotNil(t, report.Report)
		for _, r := range *report.Report {
			assert.NotNil(t, r.Statuses)
			assert.Equal(t, len(r.Statuses), 1)
			assert.Equal(t, len(r.Statuses[0].Distributions), 2)
		}
	})
}

func TestDonorReport(t *testing.T) {
	t.Run("test single item type", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		registrationData := testutils.CreateTestDonationRegistration("Andrew", 42, api.Clothing, &description)
		_, _ = dbstore.InsertRegistration(registrationData)

		distributionData := testutils.CreateTestDonationDistribution(1, 4, api.Clothing, &description)
		_, _ = dbstore.InsertDistribution(distributionData)
		_, _ = dbstore.InsertDistribution(distributionData)

		report, err := dbstore.GetDonorReport()
		assert.NoError(t, err)
		assert.NotNil(t, report.Report)
		assert.Equal(t, len(*report.Report), 1)
		r := (*report.Report)[0]
		assert.Equal(t, r.Name, "Andrew")
		assert.Equal(t, len(r.Donations), 1)
		donations := r.Donations[0]
		assert.Equal(t, donations.Type, api.Clothing)
		assert.Equal(t, donations.Quantity, 42)
		assert.Equal(t, donations.QuantityDistributed, 8)
	})

	t.Run("test two item types", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		for i, v := range []api.DonationType{api.Clothing, api.Food} {
			registrationData := testutils.CreateTestDonationRegistration("Andrew", 42, v, &description)
			_, _ = dbstore.InsertRegistration(registrationData)
			_, _ = dbstore.InsertRegistration(registrationData)
			distributionData := testutils.CreateTestDonationDistribution(i*2, 4, v, &description)
			_, _ = dbstore.InsertDistribution(distributionData)
			_, _ = dbstore.InsertDistribution(distributionData)
		}

		report, err := dbstore.GetDonorReport()
		assert.NoError(t, err)
		assert.NotNil(t, report.Report)
		assert.Equal(t, len(*report.Report), 1)
		donorReport := (*report.Report)[0]
		assert.Equal(t, donorReport.Name, "Andrew")
		for _, d := range donorReport.Donations {
			assert.Equal(t, d.Quantity, 84)
			assert.Equal(t, d.QuantityDistributed, 8)
		}
	})

	t.Run("test two donors", func(t *testing.T) {
		dbstore := TestingStoreSetup()
		defer dbstore.db.Close()
		description := "socks"
		for i, name := range []string{"Andrew", "Terry"} {
			registrationData := testutils.CreateTestDonationRegistration(name, 42, api.Money, &description)
			_, _ = dbstore.InsertRegistration(registrationData)
			distributionData := testutils.CreateTestDonationDistribution(i+1, 4, api.Money, &description)
			_, _ = dbstore.InsertDistribution(distributionData)
			_, _ = dbstore.InsertDistribution(distributionData)
		}

		report, err := dbstore.GetDonorReport()
		assert.NoError(t, err)
		assert.NotNil(t, report.Report)
		assert.Equal(t, len(*report.Report), 2)
		seen := map[string]bool{}
		for _, r := range *report.Report {
			seen[r.Name] = true
			for _, d := range r.Donations {
				assert.Equal(t, d.Quantity, 42)
				assert.Equal(t, d.QuantityDistributed, 8)
			}
		}
		assert.True(t, seen["Andrew"])
		assert.True(t, seen["Terry"])
		assert.False(t, seen["Ten"])
	})

}
