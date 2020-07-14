package goog

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"io/ioutil"
	"time"
)

func GoogleMQTTPassword(projectId, pkey string) (string, error) {
	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(pkey))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, &jwt.StandardClaims{
		Audience:  projectId,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:  time.Now().Unix(),
	})

	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GoogleClientID(projectID, region, registryID, gatewayId string) string {
	return fmt.Sprintf(
		"projects/%s/locations/%s/registries/%s/devices/%s",
		projectID, region, registryID, gatewayId,
	)
}

func GoogleTLSConfig(rootCAPath string) *tls.Config {
	certpool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(rootCAPath)
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

func GoogleMqttTopic(deviceId, messageType string) string {
	return fmt.Sprintf("/devices/%s/%s", deviceId, messageType)
}

func GoogleMqttEventsTopic(deviceId string) string {
	return GoogleMqttTopic(deviceId, "events")
}

func PublishEvent(c mqtt.Client, deviceId string, msg []byte) error {
	token := c.Publish(
		GoogleMqttEventsTopic(deviceId),
		1,
		false,
		msg,
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func GoogleMqttAttachTopic(deviceId string) string {
	return GoogleMqttTopic(deviceId, "attach")
}

func AttachDevice(c mqtt.Client, deviceId string) error {
	token := c.Publish(
		GoogleMqttAttachTopic(deviceId),
		1,
		false,
		"{}",
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func GoogleMqttDetachTopic(deviceID string) string {
	return GoogleMqttTopic(deviceID, "detach")
}

func DetachDevice(c mqtt.Client, deviceID string) error {
	token := c.Publish(
		GoogleMqttDetachTopic(deviceID),
		1,
		false,
		"{}",
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func GoogleMqttStateTopic(deviceID string) string {
	return GoogleMqttTopic(deviceID, "state")
}

func SetDeviceState(c mqtt.Client, deviceID string, state []byte) error {
	token := c.Publish(
		GoogleMqttStateTopic(deviceID),
		1,
		false,
		state,
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func GoogleMqttCommandsTopic(deviceID string) string {
	return GoogleMqttTopic(deviceID, "commands/#")
}

func WatchDeviceCommands(ctx context.Context, c mqtt.Client, deviceID string, commands chan<- []byte) error {
	token := c.Subscribe(GoogleMqttCommandsTopic(deviceID), 0, func(client mqtt.Client, m mqtt.Message) {
		select {
		case commands <- m.Payload():
		default:
			// unsubscribe from the topic if the channel is closed
			c.Unsubscribe(GoogleMqttCommandsTopic(deviceID))
			return
		}
	})

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func GoogleMqttConfigTopic(deviceID string) string {
	return GoogleMqttTopic(deviceID, "config")
}

func WatchDeviceConfig(ctx context.Context, c mqtt.Client, deviceID string, configs chan<- []byte) error {
	token := c.Subscribe(GoogleMqttConfigTopic(deviceID), 1, func(client mqtt.Client, m mqtt.Message) {
		select {
		case configs <- m.Payload():
		default:
			// unsubscribe from the topic if the channel is closed
			c.Unsubscribe(GoogleMqttConfigTopic(deviceID))
			return
		}
	})

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}
