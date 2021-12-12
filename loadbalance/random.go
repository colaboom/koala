package loadbalance

import (
	"context"
	"github.com/koala/errno"
	"github.com/koala/registry"
	"math/rand"
)

type RandomBalance struct {
	name string
}

func NewRandomBalance() LoadBalance {
	return &RandomBalance{
		name: "random",
	}
}

func (b *RandomBalance) Name() string {
	return b.name
}

func (b *RandomBalance) Select(ctx context.Context, nodes []*registry.Node) (node *registry.Node, err error) {
	nodesLen := len(nodes)
	if nodesLen == 0 {
		err = errno.NotHaveInstance
	}

	// 将节点加入到已选列表，
	defer func() {
		if node != nil {
			SetSelectedNodes(ctx, node)
		}
	}()

	// 筛掉已选过的节点，重试时避免再次选中
	nodes = filterNodes(ctx, nodes)
	if len(nodes) == 0 {
		err = errno.AllNodeFailed
		return
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
	for index, node := range nodes {
		curWeight -= node.Weight
		if curWeight < 0 {
			curIndex = index
			break
		}
	}

	if curIndex == -1 {
		err = errno.NotHaveInstance
	}

	node = nodes[curIndex]
	return
}
