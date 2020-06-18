package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

const (
	MQTTConnectionDefault = "ssl://mqtt.googleapis.com:8883"
)

type Mqtt struct {
	project  string
	region   string
	registry string
	gateway  string

	pkey   string
	rootca string
}

type MqttOptionsFunc func(m *Mqtt)

func NewMqtt(project, region, registry, gateway string, pkey, rootca string, opts ...MqttOptionsFunc) *Mqtt {
	m := &Mqtt{
		project:  project,
		region:   region,
		registry: registry,
		gateway:  gateway,
		pkey:     pkey,
		rootca:   rootca,
	}

	for _, o := range opts {
		o(m)
	}

	return m
}

func (m *Mqtt) Run(ctx context.Context) {

}

func (m *Mqtt) Connect() error {
	opts := mqtt.NewClientOptions().
		AddBroker(MQTTConnectionDefault).
		SetTLSConfig(getTLSConfig(
			m.rootca,
		)).
		SetProtocolVersion(4).
		SetClientID(getClientID(
			viper.GetString(m.project),
			viper.GetString(m.region),
			viper.GetString(m.registry),
			viper.GetString(m.gateway),
		)).
		SetUsername("unused").
		SetPassword(getJWT(
			m.project,
			m.pkey,
		)).
		SetOnConnectHandler(m.onConnect).
		SetDefaultPublishHandler(m.onMessage).
		SetConnectionLostHandler(m.onDisconnect).
		SetWriteTimeout(time.Second * 3)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (m *Mqtt) onConnect(client mqtt.Client) {

}

func (m *Mqtt) onMessage(client mqtt.Client, msg mqtt.Message) {

}

func (m *Mqtt) onDisconnect(client mqtt.Client, err error) {

}

func getTLSConfig(rootca string) *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(rootca)
	if err != nil {
		panic(err)
	}
	certpool.AppendCertsFromPEM(pemCerts)

	tlsConfig := &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}

	return tlsConfig
}

func getJWT(projectId, pkey string) string {
	bb, err := ioutil.ReadFile(pkey)
	if err != nil {
		panic(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(bb)
	if err != nil {
		panic(err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		Audience:  projectId,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		panic(err)
	}

	return tokenString
}

func getClientID(projectID, region, registryID, gatewayId string) string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/registries/%s/devices/%s",
		projectID, region, registryID, gatewayId,
	)

}
