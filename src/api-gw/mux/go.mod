module github.com/mariusmagureanu/gopex/src/api-gw/mux

go 1.16

replace (
	github.com/mariusmagureanu/gopex/pexip => ../../pexip
	github.com/mariusmagureanu/gopex/pkg/errors => ../../../pkg/errors
	github.com/mariusmagureanu/gopex/pkg/log => ../../../pkg/log
)

require (
	github.com/mariusmagureanu/gopex/pexip v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/errors v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/log v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)
