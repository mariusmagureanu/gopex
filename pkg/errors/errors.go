// Package errors is a common placeholder
// for all the custom errors in pexip monitor.
package errors

import "errors"

var (
	ErrorRoomAlreadyStarted = errors.New("room monitoring is already running")
	ErrorSSEBodyIsEmpty     = errors.New("event message is empty")
)
