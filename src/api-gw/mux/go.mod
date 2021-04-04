module bitbucket.org/kinlydev/gopex/src/api-gw/mux

go 1.16

replace (
	bitbucket.org/kinlydev/gopex/pexip => ../../pexip
	bitbucket.org/kinlydev/gopex/pkg/log => ../../../pkg/log
)

require (
	bitbucket.org/kinlydev/gopex/pexip v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/log v0.0.0-00010101000000-000000000000
	github.com/gorilla/mux v1.8.0
)
