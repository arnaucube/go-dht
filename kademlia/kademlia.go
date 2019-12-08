package kademlia

import (
	"math/bits"
	"net/rpc"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	alpha       = 3  // 'alpha', max parallalel calls
	B           = 20 // 'B', 160 bits, ID length
	KBucketSize = 20 // 'K', bucket size
)

type ListedNode struct {
	ID   ID
	Addr string
	Port string
}

type Kademlia struct {
	// N is this node data
	N        ListedNode
	KBuckets [B * 8][]ListedNode
}

func NewKademliaTable(id ID, addr, port string) *Kademlia {
	return &Kademlia{
		N: ListedNode{
			ID:   id,
			Addr: addr,
			Port: port,
		},
	}
}

func (kad Kademlia) String() string {
	r := "Node ID: " + kad.N.ID.String() + ", KBuckets:\n"
	for i, kb := range kad.KBuckets {
		if len(kb) > 0 {
			r += "	KBucket " + strconv.Itoa(i) + "\n"
			for _, e := range kb {
				r += "		" + e.ID.String() + "\n"
			}
		}
	}
	return r
}

func (kad Kademlia) KBucket(o ID) int {
	d := kad.N.ID.Distance(o)
	return kBucketByDistance(d[:])

}

func kBucketByDistance(b []byte) int {
	for i := 0; i < B; i++ {
		for a := b[i] ^ 0; a != 0; a &= a - 1 {
			return (B-1-i)*8 + (7 - bits.TrailingZeros8(bits.Reverse8(uint8(a))))
		}
	}
	return (B*8 - 1) - (B*8 - 1)
}

func (kad *Kademlia) Update(o ListedNode) {
	k := kad.KBucket(o.ID)
	kb := kad.KBuckets[k]
	if len(kb) >= KBucketSize {
		// if n.KBuckets[k] is alrady full, perform ping of the first element
		log.Info("node.KBuckets[k] already full, performing ping to node.KBuckets[0]")
		kad.PingOldNode(k, o)
		return
	}
	// check that is not already in the list
	exist, pos := existsInListedNodes(kad.KBuckets[k], o)
	if exist {
		// update position of o to the bottom
		kad.KBuckets[k] = moveToBottom(kad.KBuckets[k], pos)
		log.Info("ListedNode already exists, moved to bottom")
		return
	}
	// not exists, add it to the kBucket
	kad.KBuckets[k] = append(kad.KBuckets[k], o)
	log.Info("ListedNode not exists, added to the bottom")
	return
}

func (kad *Kademlia) PingOldNode(k int, o ListedNode) {
	// TODO when rpc layer is done
	// ping the n.KBuckets[k][0] (using goroutine)
	// if no response (timeout), delete it and add 'o'
	// n.KBuckets[k][0] = o
}

func (kad *Kademlia) CallPing(o ListedNode) error {
	client, err := rpc.DialHTTP("tcp", o.Addr+":"+o.Port)
	if err != nil {
		return err
	}
	ln := ListedNode{
		ID:   kad.N.ID,
		Addr: kad.N.Addr,
		Port: kad.N.Port,
	}
	var reply ListedNode
	err = client.Call("Node.Ping", ln, &reply)
	if err != nil {
		return err
	}
	return nil
}

func (kad *Kademlia) CallPong(o ListedNode) error {
	client, err := rpc.DialHTTP("tcp", o.Addr+":"+o.Port)
	if err != nil {
		return err
	}
	ln := ListedNode{
		ID:   kad.N.ID,
		Addr: kad.N.Addr,
		Port: kad.N.Port,
	}
	var reply bool
	err = client.Call("Node.Pong", ln, &reply)
	if err != nil {
		return err
	}

	return nil
}

func existsInListedNodes(lns []ListedNode, ln ListedNode) (bool, int) {
	for i, v := range lns {
		if v.ID.Equal(ln.ID) {
			return true, i
		}
	}
	return false, 0
}

func moveToBottom(kb []ListedNode, k int) []ListedNode {
	e := kb[k]
	kb = append(kb[:k], kb[k+1:]...)
	kb = append(kb[:], e)
	return kb
}
