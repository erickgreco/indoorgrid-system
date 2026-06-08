package gopro

import (
	"tinygo.org/x/bluetooth"
)

// * UUID's according to OpenGoPro docs (BLE Setup) "https://gopro.github.io/OpenGoPro/docs"
// * GP-XXXX is shorthand for GoPro's 128-bit UUID: b5f9XXXX-aa8d-11e3-9046-0002a5d5c51b

func mapChars(chars map[string]bluetooth.DeviceCharacteristic) *GoProChars {
	gp := func(xxxx string) string {
		return "b5f9" + xxxx + "-aa8d-11e3-9046-0002a5d5c51b"
	}

	return &GoProChars{
		WiFiAPSSID:                chars[gp("0002")],
		WiFiAPPassword:            chars[gp("0003")],
		WiFiAPPower:               chars[gp("0004")],
		WiFiAPState:               chars[gp("0005")],
		NetworkManagementCommand:  chars[gp("0091")],
		NetworkManagementResponse: chars[gp("0092")],
		Command:                   chars[gp("0072")],
		CommandResponse:           chars[gp("0073")],
		Settings:                  chars[gp("0074")],
		SettingsResponse:          chars[gp("0075")],
		Query:                     chars[gp("0076")],
		QueryResponse:             chars[gp("0077")],
	}
}
