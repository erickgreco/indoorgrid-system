package gopro

import (
	"encoding/hex"
	"fmt"
	"strings"

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
