package syserrors

import "errors"

var (
	ErrDeviceNotConnected = errors.New("error, device not connected")
)
