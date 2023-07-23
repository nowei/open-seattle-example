package api

import (
	"encoding/json"
	"log"
	"net/http"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=config.yaml ../../../api/openapi.yaml

type Server struct{}

func NewServer() ServerInterface {
	return &Server{}
}
func (s *Server) RegisterDonation(w http.ResponseWriter, r *http.Request) {
	var donation DonationRegistration
	if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("error decoding the data"))
		return
	}
	// Write information to db

	// return success
}

func (s *Server) DistributeDonation(w http.ResponseWriter, r *http.Request) {
	log.Printf("Stuffs")
	var data DonationDistribution
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("error decoding the data"))
		return
	}
	// Get donation

	// Aggregate all distributions of this donation

	// If new distribution + aggregate < original donation, add a new distribution

	// return success
}

func (s *Server) GetDonationInventoryReport(w http.ResponseWriter, r *http.Request) {

	// For each type

	// For all donations for the type
	// For all distributions for the donation

	// return result
}

func (s *Server) GetDonorReport(w http.ResponseWriter, r *http.Request) {

	// For each donor

	// Aggregate donations by type

	// Aggregate distributions by type

	// return result
}
