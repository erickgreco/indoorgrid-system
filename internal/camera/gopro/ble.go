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

func GoProServices(goPro *GoPro) (*GoProChars, error) {
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
