package gopro

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const credsFile = "data/gopro.json"

func loadCreds() (*WiFiAP, error) {
	data, err := os.ReadFile(credsFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var creds WiFiAP
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func saveCreds(creds *WiFiAP) error {
	data, err := json.MarshalIndent(creds, "", " ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(credsFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(credsFile, data, 0600)
}
