package provocation

import (
	"context"

	"github.com/secretsengine/provocation/engines/consul"
	"github.com/secretsengine/provocation/engines/password"
	"github.com/secretsengine/provocation/engines/postgresql"
	"github.com/secretsengine/provocation/engines/rabbitmq"
)

var (
	_ Engine = &consul.Engine{}
	_ Engine = &password.Engine{}
	_ Engine = &postgresql.Engine{}
	_ Engine = &rabbitmq.Engine{}
)

// Engine is a secrets engine that can provision and revoke credentials.
type Engine interface {
	Provision(ctx context.Context, namespace, name string) ([]byte, map[string][]byte, error)
	Revoke(ctx context.Context, revocation []byte) error
}
