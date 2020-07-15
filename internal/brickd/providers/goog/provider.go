package goog

import (
	"context"
	"crypto/tls"
	"github.com/brickd/agent/internal/brickd"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

const (
	MQTTConnectionDefault = "ssl://mqtt.googleapis.com:8883"
)

type Agent struct {
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

func NewAgent(L *logrus.Entry, projectId, region, registryID, gatewayID, pkeyPath, rootCAPath string) (brickd.Agent, error) {
	clientID := GoogleClientID(projectId, region, registryID, gatewayID)

	password, err := GoogleMQTTPassword(projectId, pkeyPath)
	if err != nil {
		return nil, err
	}

	tlsConf := GoogleTLSConfig(rootCAPath)

	return &Agent{
		Entry:     L,
		ClientId:  clientID,
		Password:  password,
		TLSConfig: tlsConf,
		gatewayID: gatewayID,
	}, nil
}

func NewAgentFromConfig(L *logrus.Entry) (brickd.Agent, error) {
	return NewAgent(
		L,
		viper.GetString(GCPProjectIDKey),
		viper.GetString(RegionKey),
		viper.GetString(RegistryKey),
		viper.GetString(GatewayKey),
		viper.GetString(PrivateKeyKey),
		viper.GetString(RootCAKey),
	)
}

func (m *Agent) Connect(ctx context.Context) error {
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

func (m *Agent) Disconnect(ctx context.Context, waitms uint) error {
	m.client.Disconnect(waitms)
	return nil
}

func (m *Agent) Publish(msg []byte) error {
	return PublishEvent(m.client, m.gatewayID, msg)
}

func (m *Agent) PublishAs(deviceID string, msg []byte) error {
	return PublishEvent(m.client, deviceID, msg)
}

func (m *Agent) Attach(deviceID string) error {
	return AttachDevice(m.client, deviceID)
}

func (m *Agent) Detach(deviceID string) error {
	return DetachDevice(m.client, deviceID)
}

func (m *Agent) SetState(state []byte) error {
	return SetDeviceState(m.client, m.gatewayID, state)
}

func (m *Agent) SetStateAs(deviceID string, state []byte) error {
	return SetDeviceState(m.client, deviceID, state)
}

func (m *Agent) WatchConfig(ctx context.Context) (<-chan []byte, error) {
	return WatchDeviceConfig(ctx, m.client, m.gatewayID)
}

func (m *Agent) WatchConfigAs(ctx context.Context, deviceID string) (<-chan []byte, error) {
	return WatchDeviceConfig(ctx, m.client, deviceID)
}

func (m *Agent) WatchCommands(ctx context.Context) (<-chan []byte, error) {
	return WatchDeviceCommands(ctx, m.client, m.gatewayID)
}

func (m *Agent) WatchCommandsAs(ctx context.Context, deviceID string) (<-chan []byte, error) {
	return WatchDeviceCommands(ctx, m.client, deviceID)
}
