package client

import (
	"context"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	MQTTConnectionDefault = "ssl://mqtt.googleapis.com:8883"
)

type Client struct {
	*logrus.Entry

	project  string
	region   string
	registry string
	gateway  string

	pkey   string
	rootca string

	client mqtt.Client
}

type ClientOptionsFunc func(m *Client)

func NewClient(L *logrus.Entry, project, region, registry, gateway string, pkey, rootca string, opts ...ClientOptionsFunc) *Client {
	m := &Client{
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

func (m *Client) Connect(ctx context.Context) error {
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
		SetTLSConfig(GoogleTLSConfig(
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

	m.client = mqtt.NewClient(opts)
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *Client) onConnect(client mqtt.Client) {
	m.Info("Client connected")
}

func (m *Client) onMessage(client mqtt.Client, msg mqtt.Message) {

}

func (m *Client) onDisconnect(client mqtt.Client, err error) {
	m.Info("Client disconnected")
}

func (m *Client) AttachDevice(id string) error {
	m.Debug("Attaching device: ", id)

	token := m.client.Publish(
		GoogleMqttAttachTopic(id),
		1,
		false,
		"{}",
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *Client) DetachDevice(id string) error {
	m.Debug("Detaching device: ", id)

	token := m.client.Publish(
		GoogleMqttDetachTopic(id),
		1,
		false,
		"{}",
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *Client) Publish(id string, msg []byte) error {
	m.Debug("Publishing message from device: ", id, " msg:", string(msg))

	token := m.client.Publish(
		GoogleMqttEventsTopic(id),
		1,
		false,
		msg,
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (m *Client) PublishAsGateway(msg []byte) error {
	m.Debug("Publishing message as gateway: ", m.gateway, " msg:", string(msg))

	token := m.client.Publish(
		GoogleMqttEventsTopic(m.gateway),
		1,
		false,
		msg,
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}
