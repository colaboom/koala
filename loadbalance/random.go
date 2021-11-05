package loadbalance

import (
	"context"
	"github.com/koala/errno"
	"github.com/koala/registry"
	"math/rand"
)

type RandomBalance struct {
}

func (b *RandomBalance) Name() string {
	return "random"
}

func (b *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	nodesLen := len(nodes)
	if nodesLen == 0 {
		err = errno.NotHaveInstance
	}

	var totalWeight int
	for _, val := range nodes {
		if val.Weight == 0 {
			val.Weight = DefaultNodeWeight
		}
		totalWeight += val.Weight
	}
	curWeight := rand.Intn(totalWeight)
	curIndex := -1
	return
}
