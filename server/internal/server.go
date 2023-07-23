package internal

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nowei/open-seattle-example/server/internal/api"
	"github.com/nowei/open-seattle-example/server/internal/logger"
	"github.com/nowei/open-seattle-example/server/internal/store"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest --config=api/config.yaml ../../api/openapi.yaml

type Server struct {
	db *store.DbStore
}

var log = logger.GetLogger().Sugar()

func respond(w http.ResponseWriter, statusCode int, obj any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(obj)
	if err != nil {
		log.Errorf("Failed to write response: %s", err.Error())
	}
}

type Response struct {
	Message string `json:"message"`
}

func NewServer() api.ServerInterface {
	db := store.InstantiateDbStore()
	if db == nil {
		log.Fatalf("Could not load db, exiting")
	}
	return &Server{
		db: db,
	}
}

// Registers the donation
func (s *Server) RegisterDonation(w http.ResponseWriter, r *http.Request) {
	var registrationData api.DonationRegistration
	if err := json.NewDecoder(r.Body).Decode(&registrationData); err != nil {
		respond(w, http.StatusBadRequest, Response{Message: "error decoding the data"})
		return
	}

	registration, err := s.db.InsertRegistration(registrationData)
	if err != nil {
		log.Errorf("Could not register donation: %s", err.Error())
		respond(w, http.StatusBadRequest, Response{Message: "could not register donation"})
	} else {
		respond(w, http.StatusCreated, *registration)
	}
}

// Distributes the donation
func (s *Server) DistributeDonation(w http.ResponseWriter, r *http.Request) {
	var distributionData api.DonationDistribution
	if err := json.NewDecoder(r.Body).Decode(&distributionData); err != nil {
		respond(w, http.StatusBadRequest, Response{Message: "error decoding the data"})
		return
	}
	donation, err := s.db.GetDonationRegistration(distributionData.DonationId)
	if err != nil {
		respond(w, http.StatusBadRequest, Response{Message: "error getting donation"})
		return
	}

	amountDistributed, err := s.db.GetDistributedDonationAmount(distributionData.DonationId)
	if err != nil {
		respond(w, http.StatusBadRequest, Response{Message: "error getting already donated amount"})
		return
	}
	// If new distribution + aggregate < original donation, add a new distribution
	if amountDistributed+distributionData.Quantity > donation.Quantity {
		remaining := donation.Quantity - amountDistributed
		respond(w, http.StatusBadRequest, Response{Message: fmt.Sprintf("donation distribution will exceed original donation amount, remaining: %d, input: %d", remaining, distributionData.Quantity)})
		return
	}

	distribution, err := s.db.InsertDistribution(distributionData)
	if err != nil {
		log.Errorf("Could not distribute donation: %s", err.Error())
		respond(w, http.StatusBadRequest, Response{Message: "could not distribute donation"})
	} else {
		respond(w, http.StatusCreated, *distribution)
	}
}

// Gets the donation inventory report
func (s *Server) GetDonationInventoryReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.db.GetInventoryReport()
	if err != nil {
		log.Errorf("Could not create donation inventory report: %s", err.Error())
		respond(w, http.StatusInternalServerError, Response{Message: "could not create inventory report"})
	} else {
		respond(w, http.StatusOK, *report)
	}
}

// Gets the donor report
func (s *Server) GetDonorReport(w http.ResponseWriter, r *http.Request) {
	report, err := s.db.GetDonorReport()
	if err != nil {
		log.Errorf("Could not create donor report: %s", err.Error())
		respond(w, http.StatusInternalServerError, Response{Message: "could not create donor report"})
	} else {
		respond(w, http.StatusOK, *report)
	}
}
