package gopro

import (
	"encoding/binary"
	"fmt"
	"time"

	gopropb "github.com/erickgreco/indoorgrid-system/internal/camera/gopro/proto"
	"github.com/erickgreco/indoorgrid-system/pkg/logger"
	"github.com/erickgreco/indoorgrid-system/pkg/syserrors"
	"google.golang.org/protobuf/proto"
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

	g.device = &device

	chars, err := g.GetCharacteristics()
	if err != nil {
		return logger.Error(logger.CharsServErr, err)
	}
	g.chars = *chars

	queryCh, err := g.EnableGoProNotifications(g.chars.QueryResponse)
	if err != nil {
		return logger.Error(logger.NotificationsEnableErr, err)
	}

	commandCh, err := g.EnableGoProNotifications(g.chars.CommandResponse)
	if err != nil {
		return logger.Error(logger.NotificationsEnableErr, err)
	}

	// ! when working with notifications it is better to create one channel per char
	// ! to avoid duplicated data
	g.queryRespCh = queryCh
	g.commandRespCh = commandCh

	return nil
}

func (g *GoPro) GetCharacteristics() (*GoProChars, error) {
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

// ! Deprecated, left for documentation
func (g *GoPro) GetAvailablePresets() ([]string, error) {
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

	resp, err := g.WriteQueryWithResponse(chars.Query, payload, responseCh, featureID, responseActionID)
	if err != nil {
		return nil, logger.Error(logger.GetAvailPresetsErr, err)
	}

	presets := make([]string, len(resp))
	for i, v := range resp {
		presets[i] = fmt.Sprintf("%02X", v)
	}

	return presets, nil
}

func (g *GoPro) EnableGoProNotifications(char bluetooth.DeviceCharacteristic) (<-chan []byte, error) {
	responseCh := make(chan []byte, 64)

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

// * Drain stale packets from a previous response before writing a new command
// * A labeled break is required because a plain break inside a select only exists
// * in the select, not the enclosing for loop
// * This ensures that each query starts with a clean channel
func (g *GoPro) WriteQueryWithResponse(char bluetooth.DeviceCharacteristic, payload []byte, responseCh <-chan []byte, featureID, responseActionID byte) ([]byte, error) {
drain:
	for {
		select {
		case <-responseCh:
		default:
			break drain
		}
	}

	if _, err := char.WriteWithoutResponse(payload); err != nil {
		return nil, logger.Error(logger.WriteErr, err)
	}

	resp, err := readResponse(responseCh)
	if err != nil {
		return nil, logger.Error(logger.ReadErr, err)
	}

	if len(resp) < 2 {
		return nil, logger.Error(logger.ShortResp, syserrors.ErrResponseTooShort)
	}

	if resp[0] != featureID || resp[1] != responseActionID {
		return nil, logger.Error(logger.QueryIDErr, syserrors.ErrQueryIDMatch)
	}

	return resp[2:], nil
}

func (g *GoPro) WriteCommandWithResponse(char bluetooth.DeviceCharacteristic, payload []byte, responseCh <-chan []byte, commandID byte) ([]byte, error) {
drain:
	for {
		select {
		case <-responseCh:
		default:
			break drain
		}
	}

	if _, err := char.WriteWithoutResponse(payload); err != nil {
		return nil, logger.Error(logger.WriteErr, err)
	}

	resp, err := readResponse(responseCh)
	if err != nil {
		return nil, logger.Error(logger.ReadErr, err)
	}

	if len(resp) < 2 {
		return nil, logger.Error(logger.ShortResp, syserrors.ErrResponseTooShort)
	}

	if resp[0] != commandID {
		return nil, logger.Error(logger.CmdIDErr, syserrors.ErrCmdIDMatch)
	}

	return resp[2:], nil
}

// * GoPro recommends to always use Extended (13-bit) packet headers
// * when sending messages
// * FeatureID | ActionID | QueryID | Message
// * https://gopro.github.io/OpenGoPro/docs/ble/presets/#get-available-presets
// * proto package allows to use message field and data reading
func (g *GoPro) GetPresets() ([]PresetInfo, error) {
	req := &gopropb.RequestGetPresetStatus{}

	body, err := proto.Marshal(req)
	if err != nil {
		return nil, logger.Error(logger.ErrEncodingMsg, err)
	}

	const (
		featureID        byte = 0xF5
		actionID         byte = 0x72
		responseActionID byte = 0xF2
	)

	payload := buildPacket([]byte{featureID, actionID}, body)

	resp, err := g.WriteQueryWithResponse(g.chars.Query, payload, g.queryRespCh, featureID, responseActionID)
	if err != nil {
		return nil, logger.Error(logger.GetAvailPresetsErr, err)
	}

	var status gopropb.NotifyPresetStatus
	opts := proto.UnmarshalOptions{AllowPartial: true, DiscardUnknown: true}

	if err := opts.Unmarshal(resp, &status); err != nil {
		return nil, logger.Error(logger.ErrDecodingMsg, err)
	}

	var presets []PresetInfo

	for _, group := range status.GetPresetGroupArray() {
		for _, preset := range group.GetPresetArray() {
			presets = append(presets, PresetInfo{
				ID:        preset.GetId(),
				Title:     preset.GetTitleId().String(),
				Mode:      preset.GetMode().String(),
				IsVisible: preset.GetIsVisible(),
				GroupID:   group.GetId().String(),
			})
		}
	}

	return presets, nil
}

func (g *GoPro) LoadPreset(presetID int32) ([]byte, error) {
	const commandID byte = 0x40 // * Load Preset

	presetIDBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(presetIDBytes, uint32(presetID))

	body := append([]byte{byte(len(presetIDBytes))}, presetIDBytes...)
	payload := buildPacket([]byte{commandID}, body)

	resp, err := g.WriteCommandWithResponse(g.chars.Command, payload, g.commandRespCh, commandID)
	if err != nil {
		return nil, logger.Error(logger.LoadPresetErr, err)
	}

	return resp, nil
}

func (g *GoPro) Sleep() ([]byte, error) {
	const commandID byte = 0x05 // * Sleep

	payload := buildPacket([]byte{commandID}, nil)

	resp, err := g.WriteCommandWithResponse(g.chars.Command, payload, g.commandRespCh, commandID)
	if err != nil {
		return nil, logger.Error(logger.SleepErr, err)
	}

	return resp, nil
}
