package logger

import (
	"log/slog"
	"os"
)

const (
	BleEnableFailed              = "Failed to enable bluetooth adapter"
	DeviceLocated                = "Device located"
	ScanFailed                   = "Bluetooth scan failed"
	DeviceNotFound               = "Device not found within limit time"
	DeviceNotFoundWithinAttempts = "Device not found within attempts"
	DeviceConnErr                = "Failed connecting to device"
	DeviceConn                   = "Device successfully connected"
	GetConnErr                   = "Failed to get device connection status"
	NotConn                      = "Fevice not connected"
	DiscServErr                  = "Failed to discover services"
	CharsServErr                 = "Failed to discover characteristics"
	NotificationsEnableErr       = "Failed to enable notifications"
	NotificationsEnable          = "Notifications enabled"
	Characteristic               = "Characteristic"
	WriteErr                     = "Failed to write to characteristic"
	WriteSuccess                 = "write successfully performed"
	Payload                      = "Payload"
	Response                     = "Response"
	ShortResp                    = "Response is too short"
	CmdIDErr                     = "Unexpected command id in response"
	QueryIDErr                   = "Unexpected query id in response"
	Timeout                      = "Timeout"
	ReadErr                      = "Failed to read response"
	GetAvailPresetsErr           = "Failed to get available presets"
	ErrEncodingMsg               = "Failed to marshal request"
	ErrDecodingMsg               = "Failed to unmarshal request"
	LoadPresetErr                = "Failed to load preset"
	SleepErr                     = "Failed to put camera into sleep mode"
)

var l *slog.Logger

func Init(level string) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(level)); err != nil {
		lvl = slog.LevelInfo
	}

	l = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(l)
}

func Info(msg string, args ...any)  { l.Info(msg, args...) }
func Debug(msg string, args ...any) { l.Debug(msg, args...) }
func Warn(msg string, args ...any)  { l.Warn(msg, args...) }
func Error(msg string, err error, args ...any) error {
	l.Error(msg, append([]any{"err", err}, args...)...)
	return err
}
