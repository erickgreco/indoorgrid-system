package gopro

import "tinygo.org/x/bluetooth"

type WiFiAP struct {
	SSID     string
	Password string
	Status   string
}

type GoProChars struct {
	WiFiAPSSID                bluetooth.DeviceCharacteristic // GP-0002
	WiFiAPPassword            bluetooth.DeviceCharacteristic // GP-0003
	WiFiAPPower               bluetooth.DeviceCharacteristic // GP-0004
	WiFiAPState               bluetooth.DeviceCharacteristic // GP-0005
	NetworkManagementCommand  bluetooth.DeviceCharacteristic // GP-0091
	NetworkManagementResponse bluetooth.DeviceCharacteristic // GP-0092
	Command                   bluetooth.DeviceCharacteristic // GP-0072
	CommandResponse           bluetooth.DeviceCharacteristic // GP-0073
	Settings                  bluetooth.DeviceCharacteristic // GP-0074
	SettingsResponse          bluetooth.DeviceCharacteristic // GP-0075
	Query                     bluetooth.DeviceCharacteristic // GP-0076
	QueryResponse             bluetooth.DeviceCharacteristic // GP-0077
}

type GoPro struct {
	device      *bluetooth.Device
	chars       GoProChars
	queryRespCh <-chan []byte
}

type PresetInfo struct {
	ID        int32
	Title     string
	Mode      string
	IsVisible bool
	GroupID   string
}
