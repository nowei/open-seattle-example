package main

import (
	"encoding/json"
	"net/http"

	oapichimiddleware "github.com/deepmap/oapi-codegen/pkg/chi-middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/nowei/open-seattle-example/server/internal"
	"github.com/nowei/open-seattle-example/server/internal/api"
	logger "github.com/nowei/open-seattle-example/server/internal/logger"
)

var log = logger.GetLogger().Sugar()

func main() {

	server := api.HandlerWithOptions(internal.NewServer(), api.ChiServerOptions{})

	swagger, err := api.GetSwagger()
	if err != nil {
		log.Fatalf("error loading swagger spec: %s", err.Error())
	}

	options := oapichimiddleware.Options{
		ErrorHandler: func(w http.ResponseWriter, message string, statusCode int) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			err := json.NewEncoder(w).Encode(struct {
				StatusCode int
				Message    string
			}{statusCode, message})
			if err != nil {
				log.Errorf("error writing response back: %s", err.Error())
			}
		},
	}

	r := chi.NewRouter()
	r.Use(oapichimiddleware.OapiRequestValidatorWithOptions(swagger, &options))

	r.Use(middleware.RealIP)
	r.Use(logger.LoggingMiddleware)
	r.Use(middleware.Recoverer)

	r.Mount("/", server)
	r.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("{\"msg\": \"ok\"}"))
		if err != nil {
			log.Errorf("Failed to write response: %s", err.Error())
		}
	}))

	log.Fatalf("%v", http.ListenAndServe(":3333", r).Error())
}
