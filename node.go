package wrr

// Node normally equals a server
type Node interface {
	// server's address
	Addr() string
	// name
	Name() string
}

var _ Node = &node{}

type node struct {
	// server name if sepcified, defaults to the addr
	name string
	// address
	addr string
}

type NodeOptions func(*node)

func WithNodeName(name string) NodeOptions {
	return func(n *node) {
		n.name = name
	}
}

func NewNode(address string, opts ...NodeOptions) Node {
	n := &node{
		addr: address,
		name: address,
	}

	for _, o := range opts {
		o(n)
	}
	return n
}

func (n *node) Name() string {
	return n.name
}

func (n *node) Addr() string {
	return n.addr
}
