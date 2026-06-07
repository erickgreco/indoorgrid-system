package gopro

import (
	"encoding/hex"
	"fmt"
	"log"
	"strings"
	"time"

	"tinygo.org/x/bluetooth"
)

func isGoPro(name string) bool {
	if len(name) < 6 {
		return false
	}
	return name[:6] == "GoPro "
}

func parseUUID(s string) (bluetooth.UUID, error) {
	s = strings.ReplaceAll(s, "-", "")
	b, err := hex.DecodeString(s)
	if err != nil {
		return bluetooth.UUID{}, fmt.Errorf("invalid UUID %q: %w", s, err)
	}

	var arr [16]byte
	copy(arr[:], b)

	return bluetooth.NewUUID(arr), nil
}

func BluetoothConn() (*WiFiAP, error) {
	if creds, err := loadCreds(); err == nil && creds != nil {
		log.Println("WiFi credentials loaded from file")
		return creds, nil
	}

	var adapter = bluetooth.DefaultAdapter
	if err := adapter.Enable(); err != nil {
		return nil, fmt.Errorf("could not enable bluetooth: %w", err)
	}

	var deviceAddress bluetooth.Address
	found := make(chan bluetooth.ScanResult, 1)
	scanErr := make(chan error, 1)

	log.Println("Searching GoPro...")

	go func() {
		scanErr <- adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) {
			if isGoPro(sr.LocalName()) {
				deviceAddress = sr.Address
				found <- sr
				adapter.StopScan()
			}
		})
	}()

	select {
	case <-found:
		log.Println("GoPro located")
	case err := <-scanErr:
		return nil, fmt.Errorf("error scanning bluetooth: %w", err)
	case <-time.After(15 * time.Second):
		adapter.StopScan()
		return nil, fmt.Errorf("GoPro not found in limit time")
	}

	device, err := adapter.Connect(deviceAddress, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to GoPro: %w", err)
	}

	log.Println("Connected via BLE") //! debugging

	time.Sleep(1 * time.Second) //! debugging

	services, err := device.DiscoverServices(nil)
	if err != nil {
		return nil, fmt.Errorf("unable to discover services: %w", err)
	}

	var cmdChar bluetooth.DeviceCharacteristic
	var cmdCharResp bluetooth.DeviceCharacteristic
	var ssidChar bluetooth.DeviceCharacteristic
	var passChar bluetooth.DeviceCharacteristic

	for _, service := range services {
		chars, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			continue
		}
		for _, char := range chars {
			switch char.UUID() {
			case Command:
				cmdChar = char
			case CommandResponse:
				cmdCharResp = char
			case WiFiAPSSID:
				ssidChar = char
			case WiFiAPPassword:
				passChar = char
			}
		}
	}

	if cmdChar.UUID() == (bluetooth.UUID{}) {
		return nil, fmt.Errorf("Command characteristic (GP-0072) not found")
	}
	if cmdCharResp.UUID() == (bluetooth.UUID{}) {
		return nil, fmt.Errorf("CommandResponse characteristic (GP-0073) not found")
	}
	if ssidChar.UUID() == (bluetooth.UUID{}) {
		return nil, fmt.Errorf("WiFiAPSSID characteristic (GP-0002) not found")
	}
	if passChar.UUID() == (bluetooth.UUID{}) {
		return nil, fmt.Errorf("WiFiAPPassword characteristic (GP-0003) not found")
	}
	log.Println("All required characteristics found") //! debugging

	responseCh := make(chan []byte, 1)

	if err := cmdCharResp.EnableNotifications(func(buf []byte) {
		data := make([]byte, len(buf))
		copy(data, buf)
		responseCh <- data
	}); err != nil {
		return nil, fmt.Errorf("error enabling notifications: %w", err)
	}
	log.Println("Notifications enabled on GP-0073") //! debugging

	//* Set_Wifi Enable: TLV frame [length=3, cmd=0x17, param_len=1, value=0x01]
	// if _, err := cmdChar.WriteWithoutResponse([]byte{0x03, 0x17, 0x01, 0x01}); err != nil {
	// 	return nil, fmt.Errorf("error sending Set_Wifi Enable: %w", err)
	// }
	// log.Println("Sent Set_Wifi Enable (0x17 with value 0x01)") //! debugging

	// select {
	// case resp := <-responseCh:
	// 	log.Printf("Set_Wifi response: %x", resp)
	// 	if len(resp) < 3 {
	// 		return nil, fmt.Errorf("Set_Wifi response too short: %x", resp)
	// 	}
	// 	if resp[1] != 0x17 {
	// 		return nil, fmt.Errorf("unexpected command ID in response: 0x%x", resp[1])
	// 	}
	// 	if resp[2] != 0x00 {
	// 		return nil, fmt.Errorf("Set_Wifi Enable failed with status: 0x%x", resp[2])
	// 	}
	// case <-time.After(5 * time.Second):
	// 	return nil, fmt.Errorf("timeout waiting for Set_Wifi Enable response")
	// }
	// log.Println("WiFi AP enabled, reading credentials...") //! debugging

	time.Sleep(2 * time.Second) //! debugging

	ssidBytes, err := readCharWithTimeout(ssidChar, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("GP-0002 SSID: %w", err)
	}
	ssid := strings.TrimRight(string(ssidBytes), "\x00")
	log.Printf("SSID read: %s", ssid) //! debugging

	passBytes, err := readCharWithTimeout(passChar, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("GP-0003 password: %w", err)
	}
	password := strings.TrimRight(string(passBytes), "\x00")
	log.Println("Password read successfully") //! debugging

	if ssid == "" {
		return nil, fmt.Errorf("SSID is empty after reading GP-0002")
	}

	creds := &WiFiAP{
		SSID:     ssid,
		Password: password,
	}

	if err := saveCreds(creds); err != nil {
		log.Printf("warning: unable to save wifi credentials: %v", err)
	}

	return &WiFiAP{SSID: ssid, Password: password}, nil
}

type charReadResult struct {
	data []byte
	err  error
}

func readCharWithTimeout(char bluetooth.DeviceCharacteristic, timeout time.Duration) ([]byte, error) {
	ch := make(chan charReadResult, 1)
	go func() {
		buf := make([]byte, 64)
		n, err := char.Read(buf)
		if err != nil {
			ch <- charReadResult{nil, err}
			return
		}
		ch <- charReadResult{buf[:n], nil}
	}()

	select {
	case result := <-ch:
		return result.data, result.err
	case <-time.After(timeout):
		return nil, fmt.Errorf("read timed out after %s (device may require BLE pairing: bluetoothctl pair <MAC>)", timeout)
	}
}
