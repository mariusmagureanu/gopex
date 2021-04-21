module github.com/mariusmagureanu/gopex/src

go 1.16

replace (
	github.com/mariusmagureanu/gopex/api-gw/mux => ./api-gw/mux
	github.com/mariusmagureanu/gopex/pexip => ./pexip
	github.com/mariusmagureanu/gopex/pkg/dbl => ../pkg/dbl
	github.com/mariusmagureanu/gopex/pkg/ds => ../pkg/ds
	github.com/mariusmagureanu/gopex/pkg/errors => ../pkg/errors
	github.com/mariusmagureanu/gopex/pkg/log => ../pkg/log
)

require (
	github.com/mariusmagureanu/gopex/api-gw/mux v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pexip v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/dbl v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/log v0.0.0-00010101000000-000000000000
)
