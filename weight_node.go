package wrr

// weighted node

type WeightedNode interface {
	Node
	// the weight of the server node
	Wight() float64
}

type weightedNode struct {
	Node
	// weight of the server
	weight float64
}

type WeightedNodeOpts func(*weightedNode)

func WithWeigtOpt(w float64) WeightedNodeOpts {
	return func(n *weightedNode) {
		n.weight = w
	}
}

func NewWeightedNode(node Node, opts ...WeightedNodeOpts) WeightedNode {

	wn := &weightedNode{
		Node:   node,
		weight: 1.0,
	}

	for _, o := range opts {
		o(wn)
	}
	return wn
}

func (wn *weightedNode) Wight() float64 {
	return wn.weight
}
