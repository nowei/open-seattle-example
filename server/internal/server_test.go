package internal

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nowei/open-seattle-example/server/internal/api"
	"github.com/nowei/open-seattle-example/server/internal/store"
	"github.com/nowei/open-seattle-example/server/internal/testutils"
	"github.com/stretchr/testify/assert"
)

const DBConnection = "file:test.db?mode=memory&cache=shared"
const TestSchemaFile = "store/schemas.sql"

func TestingStoreSetup() *store.DbStore {
	db, err := sql.Open("sqlite3", DBConnection)
	if err != nil {
		log.Fatalf("could not open db: %s", err.Error())
	}

	dbStore := store.InstantiateDbStore(db, TestSchemaFile)
	if dbStore == nil {
		log.Fatalf("could not create db")
	}
	return dbStore
}

func testNewServer() *Server {
	db := TestingStoreSetup()
	return &Server{
		db: db,
	}
}

func TestRegistration(t *testing.T) {
	t.Run("test add registrations", func(t *testing.T) {
		server := testNewServer()
		defer server.db.Close()
		w := httptest.NewRecorder()
		obj := testutils.CreateTestDonationRegistration("Andrew", 5, api.Clothing, nil)
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(obj)
		assert.NoError(t, err)

		r := httptest.NewRequest(http.MethodPost, "/donation/register", &b)
		server.RegisterDonation(w, r)
		var registration api.DonationRegistration
		err = json.NewDecoder(w.Result().Body).Decode(&registration)
		defer w.Result().Body.Close()
		assert.NoError(t, err)
		assert.Nil(t, registration.Description)
		assert.Equal(t, *registration.Id, 1)
		assert.Equal(t, registration.Name, "Andrew")
		assert.Equal(t, registration.Type, api.Clothing)
		assert.Equal(t, registration.Quantity, 5)
	})
}

func TestDistribution(t *testing.T) {
	t.Run("test add multiple distributions", func(t *testing.T) {
		server := testNewServer()
		defer server.db.Close()
		w := httptest.NewRecorder()
		obj := testutils.CreateTestDonationRegistration("Andrew", 5, api.Clothing, nil)
		var b bytes.Buffer
		err := json.NewEncoder(&b).Encode(obj)
		assert.NoError(t, err)
		r := httptest.NewRequest(http.MethodPost, "/donation/register", &b)
		server.RegisterDonation(w, r)

		distributionObj := testutils.CreateTestDonationDistribution(1, 1, api.Clothing, nil)
		r = httptest.NewRequest(http.MethodPost, "/donation/distribute", &b)

		for i := 1; i <= 5; i++ {
			w = httptest.NewRecorder()
			err = json.NewEncoder(&b).Encode(distributionObj)
			assert.NoError(t, err)
			server.DistributeDonation(w, r)

			var distribution api.DonationDistribution
			err = json.NewDecoder(w.Result().Body).Decode(&distribution)
			defer w.Result().Body.Close()
			if err != nil {
				log.Fatalln(err)
			}

			assert.NoError(t, err)
			assert.Nil(t, distribution.Description)
			assert.Equal(t, *distribution.Id, i)
			assert.Equal(t, distribution.DonationId, 1)
			assert.Equal(t, distribution.Type, api.Clothing)
			assert.Equal(t, distribution.Quantity, 1)
		}

		// Don't allow distributions beyond the amount given
		w = httptest.NewRecorder()
		server.DistributeDonation(w, r)
		assert.Equal(t, w.Code, 400)
	})
}
