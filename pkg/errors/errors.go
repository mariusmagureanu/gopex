package errors

import "errors"

var (
	ErrorRoomAlreadyStarted = errors.New("room monitoring is already running")
)
