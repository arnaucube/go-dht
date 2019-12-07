package main

import (
	"math/bits"
	"strconv"
)

const (
	alpha       = 3  // 'alpha', max parallalel calls
	B           = 20 // 'B', 160 bits, ID length
	KBucketSize = 20 // 'K', bucket size
)

type ListedNode struct {
	ID   ID
	Addr string
}

type Node struct {
	ID       ID
	KBuckets [B * 8][]ListedNode
}

func (n Node) String() string {
	r := "Node ID: " + n.ID.String() + ", KBuckets:\n"
	for i, kb := range n.KBuckets {
		if len(kb) > 0 {
			r += "	KBucket " + strconv.Itoa(i) + "\n"
			for _, e := range kb {
				r += "		" + e.ID.String() + "\n"
			}
		}
	}
	return r
}

func NewNode() (Node, error) {
	// TODO if node already has id, import it
	id, err := NewID()
	if err != nil {
		return Node{}, err
	}

	var n Node
	n.ID = id
	return n, nil
}

func LoadNode(idStr string) (Node, error) {
	id, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	if err != nil {
		return Node{}, err
	}
	var n Node
	n.ID = id
	return n, nil
}

func (n Node) KBucket(o ID) int {
	d := n.ID.Distance(o)
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

func (n *Node) Update(o ListedNode) {
	k := n.KBucket(o.ID)
	kb := n.KBuckets[k]
	if len(kb) >= KBucketSize {
		// TODO ping the kb[0], and if no response, delete it, and add the current o (ID)
		return
	}
	// check that is not already in the list
	exist, _ := existsInListedNodes(n.KBuckets[k], o)
	if exist {
		// update position to the bottom
		return
	}
	// not exists, add it to the kBucket
	n.KBuckets[k] = append(n.KBuckets[k], o)
	return
}

func existsInListedNodes(lns []ListedNode, ln ListedNode) (bool, int) {
	for i, v := range lns {
		if v.ID.Equal(ln.ID) {
			return true, i
		}
	}
	return false, 0
}
