package main

import (
	"context"
	"encoding/json"
	"github.com/brickd/agent/internal/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

const (
	ProjectKey     = "project"
	ProjectDefault = "brickd"

	RegionKey     = "region"
	RegionDefault = "europe-west1"

	RegistryKey     = "registry"
	RegistryDefault = ""

	GatewayKey     = "gateway"
	GatewayDefault = ""

	PrivateKeyKey     = "privateKey"
	PrivateKeyDefault = ""

	RootCAKey     = "rootkey"
	RootCADefault = "roots.pem"
)

const (
	LogLevelDefault = "debug"
	LogLevelKey     = "log.level"
)

func init() {
	viper.SetDefault(ProjectKey, ProjectDefault)
	viper.SetDefault(RegionKey, RegionDefault)
	viper.SetDefault(RegistryKey, RegistryDefault)

	viper.SetDefault(GatewayKey, GatewayDefault)
	viper.SetDefault(PrivateKeyKey, PrivateKeyDefault)
	viper.SetDefault(RootCAKey, RootCADefault)

	viper.SetDefault(LogLevelKey, LogLevelDefault)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("device")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	viper.AutomaticEnv()
}

func main() {
	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		logrus.Fatalln("Cannot parse provided log level:", err)
	}
	logrus.SetLevel(level)

	ctx := context.Background()
	L := logrus.NewEntry(logrus.StandardLogger()).WithContext(ctx)

	listedConfig, _ := json.MarshalIndent(viper.AllSettings(), "", "  ")
	L.Infof("Starting gateway with configuration: %s", listedConfig)

	c := client.NewMqtt(
		L.WithField("Component", "MQTT Client"),
		viper.GetString(ProjectKey),
		viper.GetString(RegionKey),
		viper.GetString(RegistryKey),
		viper.GetString(GatewayKey),
		viper.GetString(PrivateKeyKey),
		viper.GetString(RootCAKey),
	)

	if err = c.Run(ctx); err != nil {
		panic(err)
	}
}
