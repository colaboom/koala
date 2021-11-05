package loadbalance

import (
	"context"
	"github.com/koala/registry"
)

const (
	DefaultNodeWeight = 100
)

type LoadBalance interface {
	Name() string
	Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error)
}
