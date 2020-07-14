package goog

import (
	"context"
	"crypto/tls"
	"github.com/brickd/agent/internal/brickd"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	MQTTConnectionDefault = "ssl://mqtt.googleapis.com:8883"
)

type Conn struct {
	*logrus.Entry

	ClientId  string
	Password  string
	TLSConfig *tls.Config

	gatewayID string

	onConnect    mqtt.OnConnectHandler
	onDisconnect mqtt.ConnectionLostHandler
	onMessage    mqtt.MessageHandler

	client mqtt.Client
}

func NewConn(L *logrus.Entry, projectId, region, registryID, gatewayID, pkeyPath, rootCAPath string) (brickd.Conn, error) {
	clientID := GoogleClientID(projectId, region, registryID, gatewayID)

	password, err := GoogleMQTTPassword(projectId, pkeyPath)
	if err != nil {
		return nil, err
	}

	tlsConf := GoogleTLSConfig(rootCAPath)

	return &Conn{
		Entry:     L,
		ClientId:  clientID,
		Password:  password,
		TLSConfig: tlsConf,
		gatewayID: gatewayID,
	}, nil
}

func (m *Conn) Connect(ctx context.Context) error {
	opts := mqtt.NewClientOptions().
		AddBroker(MQTTConnectionDefault).
		SetTLSConfig(m.TLSConfig).
		SetProtocolVersion(4).
		SetClientID(m.ClientId).
		SetUsername("unused").
		SetPassword(m.Password).
		SetOnConnectHandler(m.onConnect).
		SetConnectionLostHandler(m.onDisconnect).
		SetDefaultPublishHandler(m.onMessage).
		SetWriteTimeout(time.Second * 3)

	m.client = mqtt.NewClient(opts)
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *Conn) Publish(msg []byte) error {
	return PublishEvent(m.client, m.gatewayID, msg)
}

func (m *Conn) PublishAs(deviceId string, msg []byte) error {
	return PublishEvent(m.client, deviceId, msg)
}
