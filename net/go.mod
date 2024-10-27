module github.com/tmeisel/glib/net

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/mvrilo/go-redoc v0.1.5
	github.com/stretchr/testify v1.9.0
	github.com/tmeisel/glib/error v0.0.1
	github.com/tmeisel/glib/pagination v0.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tmeisel/glib/utils v0.0.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/tmeisel/glib/error => ../error
	github.com/tmeisel/glib/pagination => ../pagination
	github.com/tmeisel/glib/utils => ../utils
)
