package brickd

import "encoding/json"

type DeviceConfig struct {
	Components []Component `json:"components"`
}

type Component struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`
}

func ParseDeviceConfig(bb []byte) (DeviceConfig, error) {
	cfg := DeviceConfig{}
	err := json.Unmarshal(bb, &cfg)

	return cfg, err
}
