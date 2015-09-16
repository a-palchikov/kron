package service

type (
	Node interface {
		IsActive() bool
	}

	NodeId uint64

	Cluster interface {
		AddNode(host string) (NodeId, error)
		Nodes() ([]Node, error)
	}
)
