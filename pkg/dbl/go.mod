module bitbucket.org/kinlydev/gopex/pkg/dbl

go 1.16

require (
	bitbucket.org/kinlydev/gopex/pkg/ds v0.0.0-00010101000000-000000000000
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.4
)

replace bitbucket.org/kinlydev/gopex/pkg/ds => ../ds
