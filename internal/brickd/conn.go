package brickd

import "context"

type Conn interface {
	Connect(ctx context.Context) error

	Publish(msg []byte) error
	PublishAs(deviceID string, msg []byte) error

	Attach(deviceID string) error
	Detach(deviceID string) error

	SetState(state []byte) error
	SetStateAs(deviceID string, state []byte) error

	WatchConfig(ctx context.Context, configs chan<- []byte) error
	WatchConfigAs(ctx context.Context, deviceID string, configs chan<- []byte) error

	WatchCommands(ctx context.Context, commands chan<- []byte) error
	WatchCommandsAs(ctx context.Context, deviceID string, commands chan<- []byte) error
}

type ConnOptionsFunc func(m *Conn)
