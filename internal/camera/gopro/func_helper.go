package gopro

import (
	"time"

	"github.com/erickgreco/indoorgrid-system/pkg/logger"
	"github.com/erickgreco/indoorgrid-system/pkg/syserrors"
)

func isGoPro(name string) bool {
	if len(name) < 6 {
		return false
	}
	return name[:6] == "GoPro "
}

/*
* Reassebles a fragmented BLE response from the GoPro into a single complete payload
* OpenGoPro packet structure:
*	- Start packet (bit 7 = 0): 2-byte header, bits [12:0] = total msg length
* 	- Continuation (bit 7 = 0): 1-byte header, bits [3:0] = sequence counter (0-15, wraps)

* BlueZ/D-Bus does not guarantee callback delivery order even when packets arrive
* in order over the air. Adjacent packets dispatched microseconds apart may arrive inverted

* Out of order continuations are buffered by counter and flushed in sequence as soon as
* the missing predecessor arrives

* Duplicate packets (same counter already in pending) are silently dropped

* Returns the reassembled payload trimmed to totalLen, or ErrTimeout after %s
* https://gopro.github.io/OpenGoPro/docs/ble/protocol/data_protocol/
 */
func readResponse(responseCh <-chan []byte) ([]byte, error) {
	var totalLen int
	var full []byte
	pending := make(map[int][]byte)
	started := false
	next := 0

	for {
		select {
		case p := <-responseCh:
			if len(p) == 0 {
				continue
			}

			if p[0]>>7 == 0 {
				if started {
					continue
				}

				headerType := (p[0] >> 5) & 0x03
				var payloadStart int

				switch headerType {
				case 0x00: // GENERAL — longitud de 5 bits
					totalLen = int(p[0] & 0x1F)
					payloadStart = 1
				case 0x01: // EXT_13 — longitud de 13 bits
					if len(p) < 2 {
						continue
					}
					totalLen = int(p[0]&0x1F)<<8 | int(p[1])
					payloadStart = 2
				case 0x02: // EXT_16 — longitud de 16 bits
					if len(p) < 3 {
						continue
					}
					totalLen = int(p[1])<<8 | int(p[2])
					payloadStart = 3
				default: // RESERVED
					return nil, logger.Error(logger.ErrDecodingMsg, syserrors.ErrReservedHeader)
				}

				full = append(full, p[payloadStart:]...)
				started = true
			} else {
				c := int(p[0] & 0x0F)
				if _, exists := pending[c]; !exists {
					pending[c] = p[1:]
				}
			}

			if !started {
				continue
			}

			for {
				payload, ok := pending[next]
				if !ok {
					break
				}
				full = append(full, payload...)
				delete(pending, next)
				next = (next + 1) % 16
			}

			if len(full) >= totalLen {
				return full[:totalLen], nil
			}

		case <-time.After(5 * time.Second):
			return nil, logger.Error(logger.Timeout, syserrors.ErrTimeout)
		}
	}
}

func buildPacket(header, body []byte) []byte {
	msgLen := len(header) + len(body)

	packet := make([]byte, 0, 2+msgLen)
	packet = append(packet,
		0x20|byte(msgLen>>8),
		byte(msgLen&0xFF),
	)
	packet = append(packet, header...)

	return append(packet, body...)
}
