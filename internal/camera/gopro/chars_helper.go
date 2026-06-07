package gopro

import (
	"log"

	"tinygo.org/x/bluetooth"
)

// * UUID's according to OpenGoPro docs (BLE Setup) "https://gopro.github.io/OpenGoPro/docs"
// * GP-XXXX is shorthand for GoPro's 128-bit UUID: b5f9XXXX-aa8d-11e3-9046-0002a5d5c51b
var (
	Command         bluetooth.UUID // GP-0072
	CommandResponse bluetooth.UUID // GP-0073
	WiFiAPSSID      bluetooth.UUID // GP-0002
	WiFiAPPassword  bluetooth.UUID // GP-0003
	WiFiAPPower     bluetooth.UUID // GP-0004
	WiFiAPState     bluetooth.UUID // GP-0005
)

func init() {
	var err error

	Command, err = parseUUID("b5f90072-aa8d-11e3-9046-0002a5d5c51b")
	if err != nil {
		log.Fatal(err)
	}

	CommandResponse, err = parseUUID("b5f90073-aa8d-11e3-9046-0002a5d5c51b")
	if err != nil {
		log.Fatal(err)
	}

	WiFiAPSSID, err = parseUUID("b5f90002-aa8d-11e3-9046-0002a5d5c51b")
	if err != nil {
		log.Fatal(err)
	}

	WiFiAPPassword, err = parseUUID("b5f90003-aa8d-11e3-9046-0002a5d5c51b")
	if err != nil {
		log.Fatal(err)
	}

	WiFiAPPower, err = parseUUID("b5f90004-aa8d-11e3-9046-0002a5d5c51b")
	if err != nil {
		log.Fatal(err)
	}

	WiFiAPState, err = parseUUID("b5f90005-aa8d-11e3-9046-0002a5d5c51b")
	if err != nil {
		log.Fatal(err)
	}
}
