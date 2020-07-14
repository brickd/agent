package goog

import "github.com/spf13/viper"

const (
	GCPProjectIDKey     = "gcp_project_id"
	GCPProjectIDDefault = "brickd"

	RegionKey     = "region"
	RegionDefault = "europe-west1"

	RegistryKey     = "registry_id"
	RegistryDefault = ""

	GatewayKey     = "id"
	GatewayDefault = ""

	PrivateKeyKey     = "private_Key"
	PrivateKeyDefault = ""

	RootCAKey     = "rootkey"
	RootCADefault = "roots.pem"
)

func init() {
	viper.SetDefault(GCPProjectIDKey, GCPProjectIDDefault)
	viper.SetDefault(RegionKey, RegionDefault)
	viper.SetDefault(RegistryKey, RegistryDefault)

	viper.SetDefault(GatewayKey, GatewayDefault)
	viper.SetDefault(PrivateKeyKey, PrivateKeyDefault)
	viper.SetDefault(RootCAKey, RootCADefault)
}
