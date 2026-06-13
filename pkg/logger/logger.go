package logger

import (
	"log/slog"
	"os"
)

const (
	BleEnableFailed              = "BLE: Failed to enable bluetooth adapter"
	DeviceLocated                = "BLE: Device located"
	ScanFailed                   = "BLE: Bluetooth scan failed"
	DeviceNotFound               = "BLE: Device not found within limit time"
	DeviceNotFoundWithinAttempts = "BLE: Device not found within attempts"
	DeviceConnErr                = "BLE: Failed connecting to device"
	DeviceConn                   = "BLE: Device successfully connected"
	GetConnErr                   = "BLE: Failed to get device connection status"
	NotConn                      = "BLE: Device not connected"
	DiscServErr                  = "BLE: Failed to discover services"
	CharsServErr                 = "BLE: Failed to discover characteristics"
	NotificationsEnableErr       = "BLE: Failed to enable notifications"
	NotificationsEnable          = "BLE: Notifications enabled"
	Characteristic               = "BLE: Characteristic"
	WriteErr                     = "BLE: Failed to write to characteristic"
	ShortResp                    = "BLE: Response is too short"
	CmdIDErr                     = "BLE: Unexpected command id in response"
	QueryIDErr                   = "BLE: Unexpected query id in response"
	ReadErr                      = "BLE: Failed to read response"
	GetAvailPresetsErr           = "BLE: Failed to get available presets"
	ErrEncodingMsg               = "BLE: Failed to marshal request"
	ErrDecodingMsg               = "BLE: Failed to unmarshal request"
	LoadPresetErr                = "BLE: Failed to load preset"
	SleepErr                     = "BLE: Failed to put camera into sleep mode"
	Timeout                      = "Timeout"
	MQTTConnErr                  = "MQTT: Failed to connect to mqtt broker"
	ParsePayloadErr              = "MQTT: Failed to parse payload"
	SaveBME680Err                = "MQTT: Failed to save BME680 reading"
	SaveBH1750Err                = "MQTT: Failed to save BH1750 reading"
	SaveSoilErr                  = "MQTT: Failed to save soil reading"
	SubscribeErr                 = "MQTT: Failed to subscribe"
	Topic                        = "Topic"
	ReadingErr                   = "SENSORS: Sensor reading out of range"
	Reading                      = "SENSORS: Reading"
	ActiveCycleErr               = "SENSORS: Failed to get active cycle"
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
