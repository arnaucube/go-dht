package node

import (
	"go-dht/config"
	"go-dht/kademlia"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"

	log "github.com/sirupsen/logrus"
)

type Admin struct {
	node Node
	disc map[kademlia.ID][]kademlia.ListedNode
}

func NewAdmin(node Node) Admin {
	return Admin{
		node: node,
	}
}

func (a *Admin) Start() error {
	// rpc server
	err := rpc.Register(a)
	if err != nil {
		return err
	}
	//
	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	//
	rpc.HandleHTTP()
	//
	http.DefaultServeMux = oldMux
	//
	listener, err := net.Listen("tcp", ":"+config.C.AdminPort)
	if err != nil {
		return err
	}

	err = http.Serve(listener, nil)
	if err != nil {
		return err
	}
	return nil
}

func (a *Admin) Find(id kademlia.ID, lns *[]kademlia.ListedNode) error {
	log.Info("[admin-rpc] FIND ", id)

	// check if id in current node
	_, err := ioutil.ReadFile(config.C.Storage + "/" + id.String())
	if err == nil {
		*lns = []kademlia.ListedNode{
			kademlia.ListedNode{
				ID:   a.node.ID(),
				Addr: config.C.Addr,
				Port: config.C.Port,
			},
		}
		log.Info("[admin-rpc] FIND found")
		return nil
	}
	log.Info("[admin-rpc] FIND not in local Node, starting NodeLookup")

	rlns, err := a.node.Kademlia().NodeLookup(id)
	if err != nil {
		log.Debug("[admin-rpc/FIND] ERROR: ", err)
		return err
	}
	*lns = rlns

	return nil
}
