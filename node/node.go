package node

import (
	"go-dht/config"
	"go-dht/kademlia"
	"net"
	"net/http"
	"net/rpc"

	log "github.com/sirupsen/logrus"
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

func (n *Node) Start() error {
	err := rpc.Register(n)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":"+config.C.Port)
	if err != nil {
		return err
	}
	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}
	// TODO ping config.C.KnownNodes
	return nil
}

// Exposed RPC calls: Ping, Store, FindNode, FindValue

func (n *Node) Ping(ln kademlia.ListedNode, thisLn *kademlia.ListedNode) error {
	log.Info("[rpc] PING from ", ln.ID)
	n.kademlia.Update(ln)
	*thisLn = kademlia.ListedNode{
		ID:   n.ID(),
		Addr: config.C.Addr,
		Port: config.C.Port,
	}
	// TODO perform PONG call to the requester (maybe ping&pong can be unified)
	return nil
}

func (n *Node) Pong(ln kademlia.ListedNode, ack *bool) error {
	log.Info("[rpc] PONG")
	n.kademlia.Update(ln)
	return nil
}

func (n *Node) Store(data []byte, ack *bool) error {
	log.Info("[rpc] STORE")

	return nil
}

func (n *Node) FindNode() {
	log.Info("[rpc] FIND_NODE")

}

func (n *Node) FindValue() {
	log.Info("[rpc] FIND_VALUE")

}
