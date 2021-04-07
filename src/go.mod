module bitbucket.org/kinlydev/gopex/src

go 1.16

replace (
	bitbucket.org/kinlydev/gopex/api-gw/mux => ./api-gw/mux
	bitbucket.org/kinlydev/gopex/pexip => ./pexip
	bitbucket.org/kinlydev/gopex/pkg/dbl => ../pkg/dbl
	bitbucket.org/kinlydev/gopex/pkg/ds => ../pkg/ds
	bitbucket.org/kinlydev/gopex/pkg/errors => ../pkg/errors
	bitbucket.org/kinlydev/gopex/pkg/log => ../pkg/log
)

require (
	bitbucket.org/kinlydev/gopex/api-gw/mux v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pexip v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/dbl v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/log v0.0.0-00010101000000-000000000000
)
