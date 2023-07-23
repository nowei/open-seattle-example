package internal

import (
	"encoding/json"
	"net/http"

	"github.com/nowei/open-seattle-example/server/internal/api"
	"github.com/nowei/open-seattle-example/server/internal/logger"
	"github.com/nowei/open-seattle-example/server/internal/store"
)

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen --config=api/config.yaml ../../api/openapi.yaml

type Server struct {
	db *store.DbStore
}

var log = logger.GetLogger().Sugar()

func respond(w http.ResponseWriter, statusCode int, obj any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(obj)
}

func NewServer() api.ServerInterface {
	return &Server{}
}
func (s *Server) RegisterDonation(w http.ResponseWriter, r *http.Request) {
	var donation api.DonationRegistration
	if err := json.NewDecoder(r.Body).Decode(&donation); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("error decoding the data"))
		return
	}
	// return success
	// respond(w, http.StatusCreated, obj)
}

func (s *Server) DistributeDonation(w http.ResponseWriter, r *http.Request) {
	log.Infof("Stuffs")
	var data api.DonationDistribution
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(400)
		w.Write([]byte("error decoding the data"))
		return
	}

	// return success
	// respond(w, http.StatusCreated, obj)
}

func (s *Server) GetDonationInventoryReport(w http.ResponseWriter, r *http.Request) {

	// return result
	// respond(w, http.StatusOK, obj)
}

func (s *Server) GetDonorReport(w http.ResponseWriter, r *http.Request) {

	// return result
	// respond(w, http.StatusOK, obj)
}
