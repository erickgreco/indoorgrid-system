package gopro

import (
	"fmt"
	"time"

	"github.com/erickgreco/indoorgrid-system/pkg/logger"
	"github.com/erickgreco/indoorgrid-system/pkg/syserrors"
	"tinygo.org/x/bluetooth"
)

func BleConn() (*GoPro, error) {
	var adapter = bluetooth.DefaultAdapter

	if err := adapter.Enable(); err != nil {
		return nil, logger.Error(logger.BleEnableFailed, err)
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
		return nil, logger.Error(logger.DeviceNotFoundWithinAttempts, fmt.Errorf("GoPro not found after 3 attempts"))
	}

	device, err := adapter.Connect(deviceAddress, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, logger.Error(logger.DeviceConnErr, err)
	}

	logger.Info(logger.DeviceConn, "device", device)

	goPro := &GoPro{
		device: &device,
	}

	return goPro, nil
}

func (g *GoPro) GoProServices(goPro *GoPro) (*GoProChars, error) {
	conn, err := goPro.device.Connected()
	if err != nil {
		return nil, logger.Error(logger.GetConnErr, err)
	}

	if !conn {
		return nil, logger.Error(logger.NotConn, syserrors.ErrDeviceNotConnected)
	}

	services, err := goPro.device.DiscoverServices(nil)
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
func (g *GoPro) EnableNotifications(goPro *GoPro) (string, error) {
	conn, err := goPro.device.Connected()
	if err != nil {
		return "", logger.Error(logger.GetConnErr, err)
	}

	if !conn {
		return "", logger.Error(logger.NotConn, syserrors.ErrDeviceNotConnected)
	}

	chars, err := g.GoProServices(goPro)
	if err != nil {
		return "", logger.Error(logger.CharsServErr, err)
	}

	responseCh := make(chan []byte, 1)

	if err := chars.QueryResponse.EnableNotifications(func(buf []byte) {
		data := make([]byte, len(buf))
		copy(data, buf)
		responseCh <- data
	}); err != nil {
		return "", logger.Error(logger.NotificationsEnableErr, err)
	}

	logger.Info(logger.NotificationsEnable, logger.Characteristic, chars.CommandResponse)

	var payload = []byte{0x20, 0x02, 0xF5, 0x72}

	if _, err := chars.Query.WriteWithoutResponse(payload); err != nil {
		return "", logger.Error(logger.WriteErr, err)
	}

	logger.Info(logger.WriteSuccess, logger.Payload, payload)

	resp, err := readResponse(responseCh)
	if err != nil {
		return "", logger.Error(logger.ReadErr, err)
	}

	if len(resp) < 2 {
		return "", logger.Error(logger.ShortResp, syserrors.ErrResponseTooShort)
	}

	if resp[0] != 0xF5 || resp[1] != 0xF2 {
		return "", logger.Error(logger.CmdIDErr, syserrors.ErrCmdIDMatch)
	}

	return fmt.Sprintf("%X", resp[2:]), nil
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
