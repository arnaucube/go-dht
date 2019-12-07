package node

import (
	"go-dht/kademlia"
)

type Node struct {
	port     string
	kademlia *kademlia.Kademlia
}

func NewNode() (Node, error) {
	id, err := kademlia.NewID()
	if err != nil {
		return Node{}, err
	}

	var n Node
	n.kademlia = kademlia.NewKademliaTable(id)
	return n, nil
}

func LoadNode(idStr string) (Node, error) {
	id, err := kademlia.IDFromString(idStr)
	if err != nil {
		return Node{}, err
	}
	var n Node
	n.kademlia = kademlia.NewKademliaTable(id)
	return n, nil
}

func (n Node) ID() kademlia.ID {
	return n.kademlia.ID
}

// Exposed RPC calls: Ping, Store, FindNode, FindValue

func (n *Node) Ping(ln kademlia.ListedNode, ack *bool) error {
	n.kademlia.Update(ln)
	return nil
}

func (n *Node) Store(data []byte, ack *bool) error {

	return nil
}

func (n *Node) FindNode() {

}

func (n *Node) FindValue() {

}
