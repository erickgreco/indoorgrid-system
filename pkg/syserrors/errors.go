package syserrors

import "errors"

var (
	ErrDeviceNotConnected = errors.New("error, device not connected")
	ErrResponseTooShort   = errors.New("error, response is too short")
	ErrCmdIDMatch         = errors.New("error, command id in response does not match")
	ErrTimeout            = errors.New("error, timeout waiting for response")
)
