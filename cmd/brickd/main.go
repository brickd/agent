package main

import (
	"context"
	"encoding/json"
	"github.com/brickd/agent/internal/brickd/httpgateway"
	"github.com/brickd/agent/internal/brickd/providers/goog"
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

	GatewayKey     = "httpgateway"
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

	viper.SetConfigName("brickd")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/brickd")
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
	L.Infof("Starting httpgateway with configuration: %s", listedConfig)

	googClient, err := goog.NewConn(
		L.WithField("Component", "MQTT Conn"),
		viper.GetString(ProjectKey),
		viper.GetString(RegionKey),
		viper.GetString(RegistryKey),
		viper.GetString(GatewayKey),
		viper.GetString(PrivateKeyKey),
		viper.GetString(RootCAKey),
	)
	if err != nil {
		L.Error("An error occured during mqtt brickd initialization: ", err)
	}

	err = googClient.Connect(ctx)
	if err != nil {
		L.Error("An error occured during mqtt brickd connection: ", err)
	}

	gateway := httpgateway.New(
		L.WithField("Component", "Gateway"),
		googClient,
	)

	if err = gateway.RunHTTP(ctx); err != nil {
		panic(err)
	}
}
