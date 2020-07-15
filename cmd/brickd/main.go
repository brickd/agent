package main

import (
	"context"
	"encoding/json"
	"github.com/brickd/agent/internal/brickd"
	"github.com/brickd/agent/internal/brickd/httpgateway"
	"github.com/brickd/agent/internal/brickd/providers/goog"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	err := brickd.LoadConfig()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	L := logrus.NewEntry(logrus.StandardLogger()).WithContext(ctx)

	listedConfig, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	L.Infof("Starting httpgateway with configuration: %s", listedConfig)

	gagent, err := goog.NewAgentFromConfig(
		L.WithField("Component", "MQTT Agent"),
	)
	if err != nil {
		L.Error("An error occured during mqtt brickd initialization: ", err)
	}

	err = gagent.Connect(ctx)
	if err != nil {
		L.Error("An error occured during mqtt brickd connection: ", err)
	}

	gateway := httpgateway.New(
		L.WithField("Component", "Gateway"),
		gagent,
	)

	if err = gateway.RunHTTP(ctx); err != nil {
		panic(err)
	}
}
