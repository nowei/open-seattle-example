package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	api "github.com/nowei/open-seattle-example/server/internal"
	logger "github.com/nowei/open-seattle-example/server/internal/logger"
)

var log = logger.GetLogger()

func main() {

	server := api.HandlerWithOptions(api.NewServer(), api.ChiServerOptions{})
	r := chi.NewRouter()

	r.Use(middleware.RealIP)
	r.Use(logger.LoggingMiddleware)
	r.Use(middleware.Recoverer)

	r.Mount("/", server)
	r.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	http.ListenAndServe(":3333", r)
}
