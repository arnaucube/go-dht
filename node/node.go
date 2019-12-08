package node

import (
	"fmt"
	"go-dht/config"
	"go-dht/kademlia"
	"net"
	"net/http"
	"net/rpc"
	"time"

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
	n.kademlia = kademlia.NewKademliaTable(id, config.C.Addr, config.C.Port)
	return n, nil
}

func LoadNode(idStr string) (Node, error) {
	id, err := kademlia.IDFromString(idStr)
	if err != nil {
		return Node{}, err
	}
	var n Node
	n.kademlia = kademlia.NewKademliaTable(id, config.C.Addr, config.C.Port)
	return n, nil
}

func (n Node) ID() kademlia.ID {
	return n.kademlia.N.ID
}

func (n *Node) Start() error {
	// rpc server
	err := rpc.Register(n)
	if err != nil {
		return err
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", ":"+config.C.Port)
	if err != nil {
		return err
	}

	go func() {
		// TMP in order to print the KBuckets of the node
		for {
			fmt.Println(n.kademlia)
			time.Sleep(5 * time.Second)
		}
	}()

	go n.pingKnownNodes(config.C.KnownNodes)

	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}
	return nil
}

func (n *Node) pingKnownNodes(lns []kademlia.ListedNode) error {
	for _, ln := range lns {
		err := n.kademlia.CallPing(ln)
		if err != nil {
			log.Warning("[pingKnownNodes]", err)
		}
	}
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
	// perform PONG call to the requester (maybe ping&pong can be unified)
	err := n.kademlia.CallPong(ln)
	if err != nil {
		log.Warning("[PONG]", err)
	}
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
