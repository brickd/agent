package brickd

import "context"

type Agent interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context, waitms uint) error

	Publish(msg []byte) error
	PublishAs(deviceID string, msg []byte) error

	Attach(deviceID string) error
	Detach(deviceID string) error

	SetState(state []byte) error
	SetStateAs(deviceID string, state []byte) error

	WatchConfig(ctx context.Context) (<-chan []byte, error)
	WatchConfigAs(ctx context.Context, deviceID string) (<-chan []byte, error)

	WatchCommands(ctx context.Context) (<-chan []byte, error)
	WatchCommandsAs(ctx context.Context, deviceID string) (<-chan []byte, error)
}

type ConnOptionsFunc func(m *Agent)
