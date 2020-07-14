package brickd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

const (
	LogLevelDefault = "debug"
	LogLevelKey     = "log.level"

	ProjectIDKey     = "project_id"
	ProjectIDDefault = ""
)

func init() {
	viper.SetDefault(LogLevelKey, LogLevelDefault)
	viper.SetDefault(ProjectIDKey, ProjectIDDefault)

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigName("brickd")
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/brickd")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
}

func LoadConfig() error {
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	level, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return err
	}
	logrus.SetLevel(level)

	return nil
}
