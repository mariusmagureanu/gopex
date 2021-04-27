// Package errors is a common placeholder
// for all the custom errors in pexip monitor.
package errors

import "errors"

var (
	// ErrRoomAlreadyStarted is a custom error thrown
	// when a given room is already being monitored.
	ErrRoomAlreadyStarted = errors.New("room monitoring is already running")

	// ErrSSEBodyIsEmpty is a custom error thrown
	// when a received SSE has an empty message.
	ErrSSEBodyIsEmpty = errors.New("event message is empty")
)
