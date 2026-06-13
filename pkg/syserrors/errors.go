package syserrors

import "errors"

var (
	ErrDeviceNotConnected = errors.New("error, device not connected")
	ErrResponseTooShort   = errors.New("error, response is too short")
	ErrCmdIDMatch         = errors.New("error, command id in response does not match")
	ErrQueryIDMatch       = errors.New("error, query id in response does not match")
	ErrTimeout            = errors.New("error, timeout waiting for response")
	ErrReservedHeader     = errors.New("error, unexpected reserved BLE header type")
	ErrInvalidTemp        = errors.New("error, temperature out of range")
	ErrInvalidHum         = errors.New("error, humidity out of range")
	ErrInvalidPress       = errors.New("error, pressure out of range")
	ErrInvalidGas         = errors.New("error, gas resistance out of range")
	ErrInvalidLux         = errors.New("error, lux out of range")
	ErrInvalidSoil        = errors.New("error, soil moisture out of range")
)
