module github.com/mariusmagureanu/gopex/pkg/dbl

go 1.16

require (
	github.com/mariusmagureanu/gopex/pkg/ds v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/errors v0.0.0-20210427195758-84eb073d95dc
	github.com/stretchr/testify v1.5.1
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	gorm.io/driver/postgres v1.0.8
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.4
)

replace (
	github.com/mariusmagureanu/gopex/pkg/ds => ../ds
	github.com/mariusmagureanu/gopex/pkg/errors => ../errors

)
