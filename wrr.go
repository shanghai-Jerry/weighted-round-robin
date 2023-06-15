package wrr

import (
	"context"
	"sync"
)

// weightedSelector ...
type weightedSelector struct {
	// Mutex protects var below
	sync.Mutex
	// currentWight for each server, defaluts to 0
	// server’s currrentWeight should add its effective wight after each pickind end,
	// p[server] += server.Wight()
	// after each round picking, the picked server's currrentWeight minus totalWeight of all the server's effective weight
	// p[server] -= totalWeight
	currentWeight map[string]float64
}

func NewWeightedSelector() Selector {
	w := &weightedSelector{
		currentWeight: make(map[string]float64),
	}
	return w
}

func (w *weightedSelector) Name() string {
	return "wrr"
}

func (w *weightedSelector) Pick(ctx context.Context, nodes []WeightedNode) (weightedNode WeightedNode, err error) {

	var selected WeightedNode
	var totalWeight float64
	// store max currentWeight of all servers
	var maxCurrentWeight float64

	w.Lock()
	defer w.Unlock()
	for _, node := range nodes {

		// server's effective weight
		effectiveWeight := node.Wight()
		totalWeight += effectiveWeight

		// update current weight of server in each round picking
		newCurrentWeight := w.currentWeight[node.Addr()] + node.Wight()
		w.currentWeight[node.Addr()] = newCurrentWeight

		// each round picked server has the maximum currrentWeight
		if selected == nil || maxCurrentWeight < newCurrentWeight {
			selected = node
			maxCurrentWeight = newCurrentWeight
		}
	}

	// this round pick is finished, the selected is picked server
	// update picked server's current weight
	w.currentWeight[selected.Addr()] -= totalWeight

	return selected, err
}
