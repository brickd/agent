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

	PrivateKeyPathKey     = "pkey"
	PrivateKeyPathDefault = "./cert/rsa_private.pem"

	RootCAKey     = "rootkey"
	RootCADefault = "./cert/roots.pem"
)

const (
	LogLevelDefault = "info"
	LogLevelKey     = "log.level"
)

func init() {
	viper.SetDefault(ProjectKey, ProjectDefault)
	viper.SetDefault(RegionKey, RegionDefault)
	viper.SetDefault(RegistryKey, RegistryDefault)

	viper.SetDefault(GatewayKey, GatewayDefault)
	viper.SetDefault(PrivateKeyPathKey, PrivateKeyPathDefault)
	viper.SetDefault(RootCAKey, RootCADefault)

	viper.SetDefault(LogLevelKey, LogLevelDefault)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
		viper.GetString(ProjectKey),
		viper.GetString(RegionKey),
		viper.GetString(RegistryKey),
		viper.GetString(GatewayKey),
		viper.GetString(PrivateKeyPathKey),
		viper.GetString(RootCAKey),
	)

	go func() {
		c.Run(ctx)
	}()

	<-ctx.Done()
}
