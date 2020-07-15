package brickd

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"path/filepath"
)

type DeviceConfig struct {
	Components []Component `json:"components"`
}

type Component struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Source  string `json:"source"`

	Runnable bool     `json:"runnable"`
	Args     []string `json:"args"`
	Env      []string `json:"env"`
}

func ParseDeviceConfig(bb []byte) (DeviceConfig, error) {
	cfg := DeviceConfig{}
	err := json.Unmarshal(bb, &cfg)

	return cfg, err
}

func (c Component) HashName() string {
	h := md5.New()
	h.Write([]byte(c.Name + c.Source + c.Version))

	return fmt.Sprintf("%x", h.Sum(nil)) + filepath.Ext(c.Source)
}
