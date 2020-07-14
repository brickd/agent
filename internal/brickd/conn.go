package brickd

import "context"

type Conn interface {
	Connect(ctx context.Context) error

	Publish(msg []byte) error
	PublishAs(deviceId string, msg []byte) error
}

type ConnOptionsFunc func(m *Conn)
