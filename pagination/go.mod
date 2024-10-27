module github.com/tmeisel/glib/pagination

go 1.21

require (
	github.com/tmeisel/glib/error v0.0.1
	github.com/tmeisel/glib/net v0.0.1
)

replace (
	github.com/tmeisel/glib/error => ../error
	github.com/tmeisel/glib/net => ../net
	github.com/tmeisel/glib/utils => ../utils
)
