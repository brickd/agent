package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

const (
	MQTTConnectionDefault = "ssl://mqtt.googleapis.com:8883"
)

type Mqtt struct {
	*logrus.Entry

	project  string
	region   string
	registry string
	gateway  string

	pkey   string
	rootca string
}

type MqttOptionsFunc func(m *Mqtt)

func NewMqtt(L *logrus.Entry, project, region, registry, gateway string, pkey, rootca string, opts ...MqttOptionsFunc) *Mqtt {
	m := &Mqtt{
		Entry:    L,
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

func (m *Mqtt) Run(ctx context.Context) error {
	if err := m.Connect(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	return nil
}

func (m *Mqtt) Connect(ctx context.Context) error {
	gPassword, err := GoogleMQTTPassword(m.project, m.pkey)
	if err != nil {
		return err
	}

	gClientID := GoogleClientID(
		m.project,
		m.region,
		m.registry,
		m.gateway,
	)

	opts := mqtt.NewClientOptions().
		AddBroker(MQTTConnectionDefault).
		SetTLSConfig(getTLSConfig(
			m.rootca,
		)).
		SetProtocolVersion(4).
		SetClientID(gClientID).
		SetUsername("unused").
		SetPassword(gPassword).
		SetOnConnectHandler(m.onConnect).
		SetDefaultPublishHandler(m.onMessage).
		SetConnectionLostHandler(m.onDisconnect).
		SetWriteTimeout(time.Second * 3)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
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
