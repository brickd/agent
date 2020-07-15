package main

import (
	"context"
	"encoding/json"
	"github.com/brickd/agent/internal/brickd"
	"github.com/brickd/agent/internal/brickd/providers/goog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// BrickDFolder is used to store components and configurations of the agent
	BrickDFolder string

	L = logrus.NewEntry(logrus.StandardLogger())
)

func init() {
	home, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	BrickDFolder = filepath.Join(home, ".brickd")

	err = os.MkdirAll(BrickDFolder, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func main() {
	err := brickd.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	listedConfig, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	L.Infof("Starting httpgateway with configuration: %s", listedConfig)

	googClient, err := goog.NewAgentFromConfig(
		L.WithField("Component", "MQTT Agent"),
	)
	if err != nil {
		L.Error("An error occured during mqtt brickd initialization: ", err)
	}

	err = googClient.Connect(ctx)
	if err != nil {
		L.Error("An error occured during mqtt brickd connection: ", err)
	}

	watchContext, watchCancel := context.WithCancel(ctx)
	configs, err := googClient.WatchConfig(watchContext)
	if err != nil {
		L.Error("Failed to read config of the device: ", err)
	}

	cfg, err := brickd.ParseDeviceConfig(<-configs)
	if err != nil {
		L.Error("Config could not be parsed: ", err)
	}
	watchCancel()

	runContext, _ := context.WithCancel(context.Background())

	for _, c := range cfg.Components {
		L.Info("Checking component ", c.Name)
		err = CheckComponent(runContext, c)
		if err != nil {
			L.Error("Component check failed with: ", err)
		}
	}

	for _, c := range cfg.Components {
		RunComponent(runContext, c)
	}
}

func CheckComponent(ctx context.Context, c brickd.Component) error {
	if !fileExists(filepath.Join(BrickDFolder, c.Name)) {
		L.Info("Component ", c.Name, " does not exist, downloading")
		err := DownloadComponent(ctx, c)
		if err != nil {
			return err
		}
	}

	L.Info("Component ", c.Name, " checked")

	return nil
}

func RunComponent(ctx context.Context, c brickd.Component) <-chan error {
	return nil
}

func DownloadComponent(ctx context.Context, c brickd.Component) error {
	resp, err := http.Get(c.Source)
	if err != nil {
		return err
	}

	bb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		filepath.Join(BrickDFolder, c.Hash()),
		bb,
		os.ModePerm,
	)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
