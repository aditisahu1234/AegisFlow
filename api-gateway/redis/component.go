package redis

import (
	"api-gateway/config"
	"api-gateway/internal/runtime"
	"context"
)

type Component struct {
	cfg config.Config
}

func NewComponent(cfg config.Config) runtime.Component {
	return &Component{
		cfg: cfg,
	}
}

func (c *Component) Name() string {
	return "redis"
}

func (c *Component) Start(ctx context.Context) error {

	go ConnectWithRetry(
		ctx,
		c.cfg,
	)

	StartHealthMonitor()

	return nil
}

func (c *Component) Stop(ctx context.Context) error {

	if Client != nil {
		return Client.Close()
	}

	return nil
}

func (c *Component) Health(ctx context.Context) error {

	if Client == nil {
		return context.Canceled
	}

	return Client.Ping(ctx).Err()
}
