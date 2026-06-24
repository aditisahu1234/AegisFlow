package redis

import (
	"api-gateway/config"
	"api-gateway/internal/runtime"
	"context"
	"errors"
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

	return ConnectRedis(c.cfg)
}

func (c *Component) Stop(ctx context.Context) error {

	if Client != nil {
		return Client.Close()
	}

	return nil
}

func (c *Component) Health(ctx context.Context) error {

	if Client == nil {
		return errors.New("redis client not initialized")
	}

	return Client.Ping(ctx).Err()
}

func (c *Component) Dependencies() []runtime.Dependency {
	return nil
}
