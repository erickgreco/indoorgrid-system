package gopro

import (
	"fmt"
	"time"

	"github.com/erickgreco/indoorgrid-system/pkg/logger"
	"github.com/erickgreco/indoorgrid-system/pkg/syserrors"
	"tinygo.org/x/bluetooth"
)

func New() *GoPro {
	return &GoPro{}
}

func (g *GoPro) BleConn() error {
	var adapter = bluetooth.DefaultAdapter

	if err := adapter.Enable(); err != nil {
		return logger.Error(logger.BleEnableFailed, err)
	}

	var deviceAddress bluetooth.Address
	foundCh := make(chan bluetooth.ScanResult, 1)
	scanErr := make(chan error, 1)
	scanTimeOut := 5 * time.Second
	var found bool

	for i := range 3 {
		go func() {
			scanErr <- adapter.Scan(func(a *bluetooth.Adapter, sr bluetooth.ScanResult) {
				if isGoPro(sr.LocalName()) {
					deviceAddress = sr.Address
					foundCh <- sr
					adapter.StopScan()
				}
			})
		}()

		select {
		case <-foundCh:
			logger.Info(logger.DeviceLocated, "addr", deviceAddress)
			found = true
		case err := <-scanErr:
			logger.Error(logger.ScanFailed, err, "attempt", i+1)
		case <-time.After(scanTimeOut):
			adapter.StopScan()
			logger.Warn(logger.DeviceNotFound, "attempt", i+1)
		}

		if found {
			break
		}
	}
	if !found {
		return logger.Error(logger.DeviceNotFoundWithinAttempts, fmt.Errorf("GoPro not found after 3 attempts"))
	}

	device, err := adapter.Connect(deviceAddress, bluetooth.ConnectionParams{})
	if err != nil {
		return logger.Error(logger.DeviceConnErr, err)
	}

	logger.Info(logger.DeviceConn, "device", device)

	g.device = &device

	return nil
}

func (g *GoPro) GetCharacteristics() (*GoProChars, error) {
	conn, err := g.device.Connected()
	if err != nil {
		return nil, logger.Error(logger.GetConnErr, err)
	}

	if !conn {
		return nil, logger.Error(logger.NotConn, syserrors.ErrDeviceNotConnected)
	}

	services, err := g.device.DiscoverServices(nil)
	if err != nil {
		return nil, logger.Error(logger.DiscServErr, err)
	}

	discoveredChars := make(map[string]bluetooth.DeviceCharacteristic)

	for _, service := range services {
		chars, err := service.DiscoverCharacteristics(nil)
		if err != nil {
			continue
		}

		for _, char := range chars {
			discoveredChars[char.UUID().String()] = char
		}
	}

	namedChars := mapChars(discoveredChars)
	return namedChars, nil
}

// * GoPro recommends to always use Extended (13-bit) packet headers
// * when sending messages
// * FeatureID | ActionID | QueryID | Message
// * https://gopro.github.io/OpenGoPro/docs/ble/presets/#get-available-presets
// ! This function currently does not support message field
func (g *GoPro) GetAvailablePresets() ([]string, error) {
	conn, err := g.device.Connected()
	if err != nil {
		return nil, logger.Error(logger.GetConnErr, err)
	}

	if !conn {
		return nil, logger.Error(logger.NotConn, syserrors.ErrDeviceNotConnected)
	}

	chars, err := g.GetCharacteristics()
	if err != nil {
		return nil, logger.Error(logger.CharsServErr, err)
	}

	responseCh, err := g.EnableGoProNotifications(chars.QueryResponse)
	if err != nil {
		return nil, logger.Error(logger.NotificationsEnableErr, err)
	}

	var featureID byte = 0xF5
	var actionID byte = 0x72
	var responseActionID byte = 0xF2

	var payload = []byte{0x20, 0x02, featureID, actionID}

	resp, err := g.WriteWithResponse(chars.Query, payload, responseCh, featureID, responseActionID)
	if err != nil {
		return nil, logger.Error(logger.GetAvailPresetsErr, err)
	}

	presets := make([]string, len(resp))
	for i, v := range resp {
		presets[i] = fmt.Sprintf("%02X", v)
	}

	return presets, nil
}

// * GoPro documentation link: https://gopro.github.io/OpenGoPro/docs/ble/protocol/data_protocol#continuation-packets
// * When receiving a message that is longer than 20 bytes, the message
// * must be split into N packets with packet 1 containing a start packet
// * header and packets 2..N containing a continuation packet header
// * Bit 7 = 1 - continuation
// * Bit 7 = 0 - start
func readResponse(responseCh <-chan []byte) ([]byte, error) {
	var buf []byte
	var totalLen int

	for {
		packet := <-responseCh
		if len(packet) == 0 {
			continue
		}

		if packet[0]>>7 == 0 {
			totalLen = int(packet[0]&0x1F)<<8 | int(packet[1])
			buf = append(buf, packet[2:]...)
		} else {
			buf = append(buf, packet[1:]...)
		}

		if len(buf) >= totalLen {
			return buf[:totalLen], nil
		}
	}
}

func (g *GoPro) EnableGoProNotifications(char bluetooth.DeviceCharacteristic) (<-chan []byte, error) {
	responseCh := make(chan []byte, 16)

	if err := char.EnableNotifications(func(buf []byte) {
		data := make([]byte, len(buf))
		copy(data, buf)
		responseCh <- data
	}); err != nil {
		return nil, logger.Error(logger.WriteErr, err)
	}
	logger.Info(logger.NotificationsEnable, logger.Characteristic, char)

	return responseCh, nil
}

func (g *GoPro) WriteWithResponse(char bluetooth.DeviceCharacteristic, payload []byte, responseCh <-chan []byte, featureID, responseActionID byte) ([]byte, error) {
	if _, err := char.WriteWithoutResponse(payload); err != nil {
		return nil, logger.Error(logger.WriteErr, err)
	}

	logger.Info(logger.WriteSuccess, logger.Payload, payload)

	resp, err := readResponse(responseCh)
	if err != nil {
		return nil, logger.Error(logger.ReadErr, err)
	}

	if len(resp) < 2 {
		return nil, logger.Error(logger.ShortResp, syserrors.ErrResponseTooShort)
	}

	if resp[0] != featureID || resp[1] != responseActionID {
		return nil, logger.Error(logger.CmdIDErr, syserrors.ErrCmdIDMatch)
	}

	return resp, nil
}
