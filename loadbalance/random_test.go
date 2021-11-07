package loadbalance

import (
	"context"
	"fmt"
	"github.com/koala/registry"
	"testing"
)

func TestRandomBalance_Select(t *testing.T) {
	balance := &RandomBalance{}
	var weights = [3]int{50, 100, 150}
	var nodes []*registry.Node
	for i := 0; i < 3; i++ {
		node := &registry.Node{
			IP:   fmt.Sprintf("127.0.0.%d", i),
			Port: 8801,
			Weight: weights[i%3],
		}
		nodes = append(nodes, node)
	}

	var countStat = make(map[string]int)
	for i := 0; i < 1000; i++ {
		node, err := balance.Select(context.TODO(), nodes)
		if err != nil {
			t.Fatalf("Select failed, err :%v", err)
		}
		countStat[node.IP]++
	}

	for key, val := range countStat {
		fmt.Printf("ip:%s,cnt:%d\n", key, val)
	}
}
