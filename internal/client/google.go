package client

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
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

func GoogleMqttTopic(deviceId, messageType string) string {
	return fmt.Sprintf("/devices/%s/%s", deviceId, messageType)
}

func GoogleMqttEventsTopic(deviceId string) string {
	return fmt.Sprintf("/devices/%s/%s/events", deviceId, "events")
}

func GoogleMqttAttachTopic(deviceId string) string {
	return fmt.Sprintf("/devices/%s/%s/events", deviceId, "attach")
}

func GoogleMqttDetachTopic(deviceId string) string {
	return fmt.Sprintf("/devices/%s/%s/events", deviceId, "detach")
}

func GoogleMqttCommandsTopic(deviceId string) string {
	return fmt.Sprintf("/devices/%s/%s/commands/#", deviceId, "commands")
}

func GoogleMqttConfigTopic(deviceId string) string {
	return fmt.Sprintf("/devices/%s/%s/events", deviceId, "config")
}
