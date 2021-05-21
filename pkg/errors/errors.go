// Package errors is a common placeholder
// for all the custom errors in pexip monitor.
package errors

import "errors"

var (
	// ErrSSEBodyIsEmpty is a custom error thrown
	// when a received SSE has an empty message.
	ErrSSEBodyIsEmpty = errors.New("event message is empty")

	// ErrRecordNotFound is a custom error thrown when
	// a record is not found in the database. The gorm orm has
	// a similar error, the reason we're using this one is we don't
	// want clients to be dependent on the gorm package.
	ErrRecordNotFound = errors.New("record not found")

	// ErrExpiredPexipToken is a custom error thrown when an existing
	// token has expired.
	ErrExpiredPexipToken = errors.New("expired pexip token")

	// ErrNoPexipToken is a custom error thrown when no pexip token
	// has been found for the given room.
	ErrNoPexipToken = errors.New("no pexip token found for the given room")
)
