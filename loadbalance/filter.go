package loadbalance

import (
	"context"
	"fmt"
	"github.com/koala/logs"
	"github.com/koala/registry"
)

type SelectedNodes struct {
	SelectedNodeMap map[string]bool
}

type loadbalanceFilterNodes struct{}

func WithBalanceContext(ctx context.Context) context.Context {
	sel := &SelectedNodes{
		SelectedNodeMap: map[string]bool{},
	}

	return context.WithValue(ctx, loadbalanceFilterNodes{}, sel)
}

func GetSelectedNodes(ctx context.Context) *SelectedNodes {
	selectedNodes, ok := ctx.Value(loadbalanceFilterNodes{}).(*SelectedNodes)
	if !ok {
		return nil
	}
	return selectedNodes
}

func SetSelectedNodes(ctx context.Context, node *registry.Node) {
	selectedNodes := GetSelectedNodes(ctx)
	if selectedNodes == nil {
		return
	}
	addr := fmt.Sprintf("%s:%d", node.IP, node.Port)
	logs.Info(ctx, "set selected nodes : %s", addr)
	selectedNodes.SelectedNodeMap[addr] = true
}

func filterNodes(ctx context.Context, nodes []*registry.Node) (newNodes []*registry.Node) {
	selectedNodes := GetSelectedNodes(ctx)
	if selectedNodes == nil {
		return
	}
	for _, node := range nodes {
		addr := fmt.Sprintf("%s:%d", node.IP, node.Port)
		_, ok := selectedNodes.SelectedNodeMap[addr]
		if ok {
			logs.Info(ctx, "addr:%s 已经被选过", addr)
			continue
		}
		newNodes = append(newNodes, node)
	}
	return
}
