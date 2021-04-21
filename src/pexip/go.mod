module github.com/mariusmagureanu/gopex/src/pexip

replace (
	github.com/mariusmagureanu/gopex/pkg/errors => ../../pkg/errors
	github.com/mariusmagureanu/gopex/pkg/log => ../../pkg/log
)

go 1.16

require (
	github.com/mariusmagureanu/gopex/pkg/errors v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/log v0.0.0-00010101000000-000000000000
)
