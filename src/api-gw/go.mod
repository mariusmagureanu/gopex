module bitbucket.org/kinlydev/gopex/src/api-gw

go 1.16

replace (
	bitbucket.org/kinlydev/gopex/api-gw/mux => ./mux
	bitbucket.org/kinlydev/gopex/pexip => ../pexip
	bitbucket.org/kinlydev/gopex/pkg/dbl => ../../pkg/dbl
	bitbucket.org/kinlydev/gopex/pkg/ds => ../../pkg/ds
	bitbucket.org/kinlydev/gopex/pkg/log => ../../pkg/log
	bitbucket.org/kinlydev/gopex/pkg/errors => ../../pkg/errors
)

require (
	bitbucket.org/kinlydev/gopex/api-gw/mux v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pexip v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/dbl v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/log v0.0.0-00010101000000-000000000000
	golang.org/x/tools v0.1.0 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
)
