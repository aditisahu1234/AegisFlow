/*
Every infrastructure dependency (Redis, PostgreSQL,
Kafka, AI service) will implement this interface.
The Runtime Manager doesn't need to know what a
component is—it only knows how to start it, stop
it, and check its health.
*/
package runtime

import "context"

type Component interface {
	Name() string

	Start(ctx context.Context) error

	Stop(ctx context.Context) error

	Health(ctx context.Context) error

	Dependencies() []Dependency
}
