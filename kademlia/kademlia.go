package kademlia

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

type Kademlia struct {
	ID       ID
	KBuckets [B * 8][]ListedNode
}

func NewKademliaTable(id ID) *Kademlia {
	return &Kademlia{
		ID: id,
	}
}

func (kad Kademlia) String() string {
	r := "Node ID: " + kad.ID.String() + ", KBuckets:\n"
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
	d := kad.ID.Distance(o)
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
		kad.Ping(k, o)
		return
	}
	// check that is not already in the list
	exist, pos := existsInListedNodes(kad.KBuckets[k], o)
	if exist {
		// update position of o to the bottom
		kad.KBuckets[k] = moveToBottom(kad.KBuckets[k], pos)
		return
	}
	// not exists, add it to the kBucket
	kad.KBuckets[k] = append(kad.KBuckets[k], o)
	return
}

func (kad *Kademlia) Ping(k int, o ListedNode) {
	// TODO when rpc layer is done
	// ping the n.KBuckets[k][0] (using goroutine)
	// if no response (timeout), delete it and add 'o'
	// n.KBuckets[k][0] = o
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
