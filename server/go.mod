module github.com/nowei/open-seattle-example/server

go 1.20

replace github.com/nowei/open-seattle-example-server/server => ./

require (
	github.com/deepmap/oapi-codegen v1.13.0
	github.com/getkin/kin-openapi v0.118.0
	github.com/go-chi/chi/v5 v5.0.10
	github.com/mattn/go-sqlite3 v1.14.17
	go.uber.org/zap v1.24.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

require (
	github.com/go-openapi/jsonpointer v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/invopop/yaml v0.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/perimeterx/marshmallow v1.1.5 // indirect
	github.com/rogpeppe/go-internal v1.11.0 // indirect
	github.com/stretchr/testify v1.8.4
	github.com/urfave/negroni v1.0.0
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
