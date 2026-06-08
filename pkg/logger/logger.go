package logger

import (
	"log/slog"
	"os"
)

const (
	BleEnableFailed              = "failed to enable bluetooth adapter"
	DeviceLocated                = "device located"
	ScanFailed                   = "bluetooth scan failed"
	DeviceNotFound               = "device not found within limit time"
	DeviceNotFoundWithinAttempts = "device not found within attempts"
	DeviceConnErr                = "error connecting to device"
	DeviceConn                   = "device successfully connected"
	GetConnErr                   = "failed to get device connection status"
	NotConn                      = "device not connected"
	DiscServErr                  = "failed to discover services"
	CharsServErr                 = "failed to discover characteristics"
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
