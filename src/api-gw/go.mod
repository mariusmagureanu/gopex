module bitbucket.org/kinlydev/gopex/src/api-gw

go 1.16

replace (
	bitbucket.org/kinlydev/gopex/api-gw/mux => ./mux
	bitbucket.org/kinlydev/gopex/pkg/dbl => ../../pkg/dbl
	bitbucket.org/kinlydev/gopex/pkg/ds => ../../pkg/ds
	bitbucket.org/kinlydev/gopex/pkg/log => ../../pkg/log
)

require (
	bitbucket.org/kinlydev/gopex/api-gw/mux v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/dbl v0.0.0-00010101000000-000000000000
	bitbucket.org/kinlydev/gopex/pkg/ds v0.0.0-00010101000000-000000000000 // indirect
	bitbucket.org/kinlydev/gopex/pkg/log v0.0.0-00010101000000-000000000000
	gorm.io/driver/postgres v1.0.8 // indirect
	gorm.io/driver/sqlite v1.1.4 // indirect
)
