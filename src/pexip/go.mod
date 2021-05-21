module github.com/mariusmagureanu/gopex/src/pexip

replace (
	github.com/mariusmagureanu/gopex/pkg/errors => ../../pkg/errors
	github.com/mariusmagureanu/gopex/pkg/log => ../../pkg/log
)

go 1.16

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/mariusmagureanu/gopex/pkg/errors v0.0.0-00010101000000-000000000000
	github.com/mariusmagureanu/gopex/pkg/log v0.0.0-00010101000000-000000000000
	github.com/nats-io/nats-server/v2 v2.2.3 // indirect
	github.com/nats-io/nats.go v1.11.0

)
