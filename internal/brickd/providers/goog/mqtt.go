package goog

import (
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

func GoogleAttachDevice(c mqtt.Client, deviceId string) error {
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

func GoogleMqttDetachTopic(deviceId string) string {
	return GoogleMqttTopic(deviceId, "detach")
}

func GoogleDetachDevice(c mqtt.Client, deviceId string) error {
	token := c.Publish(
		GoogleMqttDetachTopic(deviceId),
		1,
		false,
		"{}",
	)

	if token.WaitTimeout(time.Second*5) == false || token.Error() != nil {
		return token.Error()
	}

	return nil
}

func GoogleMqttCommandsTopic(deviceId string) string {
	return GoogleMqttTopic(deviceId, "commands/#")
}

func GoogleMqttConfigTopic(deviceId string) string {
	return GoogleMqttTopic(deviceId, "config")
}
