package wrr

import (
	"context"
	"math/rand"
	"sort"
)

// Selector  selects server with a specfic strategy
type Selector interface {
	Name() string
	Pick(context.Context, []WeightedNode) (WeightedNode, error)
}

// accourding probablity for server selection
type defaultSelector struct{}

func NewDefaultSelector() Selector {
	return &defaultSelector{}
}

func (s *defaultSelector) Name() string {
	return "default"
}

func (d *defaultSelector) Pick(ctx context.Context, nodes []WeightedNode) (WeightedNode, error) {

	// sorting nodes according to their weight ascending
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].Wight() < nodes[j].Wight()
	})

	var totalWeight float64
	for _, node := range nodes {
		totalWeight += node.Wight()
	}

	var currTotal float64
	var selected WeightedNode
	// select server with proability
	pro := rand.Float64()
	for _, node := range nodes {
		currTotal += node.Wight()
		if pro <= currTotal/totalWeight {
			selected = node
			break
		}
	}
	return selected, nil
}
