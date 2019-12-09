package node

import (
	"encoding/hex"
	"fmt"
	"go-dht/config"
	"go-dht/kademlia"
	"io/ioutil"
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

	h := kademlia.HashData(data)
	if n.kademlia.KBucket(h) != 0 {
		*ack = false
		log.Warning("[STORE] data not for this node")
		return nil
	}
	err := ioutil.WriteFile(config.C.Storage+"/"+hex.EncodeToString(h[:]), data, 0644)
	if err != nil {
		*ack = false
		log.Warning("[STORE]", err)
		return err
	}
	*ack = true
	return nil
}

func (n *Node) FindNode(ln kademlia.ListedNode, lns *[]kademlia.ListedNode) error {
	log.Info("[rpc] FIND_NODE")
	// k := n.kademlia.KBucket(ln.ID)
	k, err := n.kademlia.FindClosestKBucket(ln.ID)
	if err != nil {
		*lns = []kademlia.ListedNode{}
		return nil
	}
	log.Info("[FIND_NODE] k", k)
	bucket := n.kademlia.KBuckets[k]
	*lns = bucket
	return nil
}

type FindValueResp struct {
	Value []byte
	Lns   []kademlia.ListedNode
}

func (n *Node) FindValue(id kademlia.ID, resp *FindValueResp) error {
	log.Info("[rpc] FIND_VALUE")
	// first check if value is in this node storage
	f, err := ioutil.ReadFile(config.C.Storage + "/" + id.String())
	if err == nil {
		// value exists, return it
		*resp = FindValueResp{
			Value: f,
		}
		return nil
	}

	// k := n.kademlia.KBucket(id)
	k, err := n.kademlia.FindClosestKBucket(id)
	if err != nil {
		*resp = FindValueResp{}
		return nil
	}
	log.Info("[FIND_VALUE] k", k)
	// bucket := n.kademlia.KBuckets[k]
	*resp = FindValueResp{
		Lns: n.kademlia.KBuckets[k],
	}
	return nil

}
