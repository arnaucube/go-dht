package kademlia

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var debug = false

func init() {
	log.SetLevel(log.DebugLevel)
}

func TestCountZeros(t *testing.T) {
	zeroes := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	assert.Equal(t, 0, kBucketByDistance(zeroes))

	b := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	assert.Equal(t, 20, len(b))
	assert.Equal(t, 7, kBucketByDistance(b))
	b[19] = 0x00
	assert.Equal(t, 7, kBucketByDistance(b))

	b[0] = 0x0f
	assert.Equal(t, 7, kBucketByDistance(b))

	b[0] = 0x0c
	assert.Equal(t, 7, kBucketByDistance(b))

	b[0] = 0x00
	b[1] = 0x00
	b[2] = 0x0f
	assert.Equal(t, 7, kBucketByDistance(b))

	b[2] = 0x07
	assert.Equal(t, 7, kBucketByDistance(b))

	b[2] = 0x03
	assert.Equal(t, 7, kBucketByDistance(b))

	b = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	b[19] = 0x01
	assert.Equal(t, 2, kBucketByDistance(b))
	b[19] = 0x05
	assert.Equal(t, 2, kBucketByDistance(b))

	b[19] = 0x10
	assert.Equal(t, 1, kBucketByDistance(b))

	b[18] = 0x10
	assert.Equal(t, 3, kBucketByDistance(b))
}

func TestKBucket(t *testing.T) {
	idA, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)
	kademlia := NewKademliaTable(idA, "127.0.0.1", "5000")

	d := kademlia.KBucket(kademlia.N.ID)
	assert.Equal(t, 0, d) // same node should have distance 0

	idB, err := IDFromString("c48d8b53dbefb609ed4e94d386dd5b22efcb2c5b")
	assert.Nil(t, err)

	d = kademlia.KBucket(idB)
	assert.Equal(t, 7, d)
}

func prepareTestListedNodes() []ListedNode {
	lnIDs := []string{
		"c48d8b53dbefb609ed4e94d386dd5b22efcb2c5b",
		"420bfebd44fc62615253224328f40f29c9b225fa",
		"6272bb67b1661abfa07d1d32cd9b810e531d42cd",
		"07db608db033384895e48098a1bbe25266387463",
		"c19c3285ab9ada4b420050ae1a204640b2bed508",
		"f8971d5da24517c8cc5a316fb0658de8906c4155",
		"04122a5f2dec42284147b1847ec1bc41ecd78626",
		"407a90870d7b482a271446c85fda940ce78a4c7a",
		"5ebe4539e7a33721a8623f7ebfab62600aa503e7",
		"fc44a42879ef3a74d6bd8303bc3e4e205a92acf9",
		"fc44a42879ef3a74d6bd8303bc3e4e205a92acf9",
		"10c86f96ebaa1685a46a6417e6faa2ef34a68976",
		"243c81515196a5b0e2d4675e73f0da3eead12793",
		"0fd85ddddf15aeec2d5d8b01b013dbca030a18d7",
		"0fd85ddddf15aeec2d5d8b01b013dbca030a18d5",
		"0fd85ddddf15aeec2d5d8b01b013dbca030a18d0",
		"0fd85ddddf15aeec2d5d8b01b013dbca030a1800",
		"0750931c40a52a2facd220a02851f7d34f95e1fa",
		"ebfba615ac50bcd0f5c2420741d032e885abf576",
	}
	var lns []ListedNode
	for i := 0; i < len(lnIDs); i++ {
		idI, err := IDFromString(lnIDs[i])
		if err != nil {
			panic(err)
		}
		lnI := ListedNode{
			ID:   idI,
			Addr: "",
		}
		lns = append(lns, lnI)
	}
	return lns
}

func TestMoveToBottom(t *testing.T) {
	lns := prepareTestListedNodes()
	movedElem := lns[3]
	assert.NotEqual(t, movedElem, lns[len(lns)-1])
	lns = moveToBottom(lns, 3)
	assert.Equal(t, movedElem, lns[len(lns)-1])
}

func TestUpdate(t *testing.T) {
	idA, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)
	kademlia := NewKademliaTable(idA, "127.0.0.1", "5000")

	lns := prepareTestListedNodes()
	for _, lnI := range lns {
		kademlia.Update(lnI)
	}

	if debug {
		fmt.Println(kademlia)
	}

	assert.Equal(t, 2, len(kademlia.KBuckets[0]))
	assert.Equal(t, 0, len(kademlia.KBuckets[1]))
	assert.Equal(t, 2, len(kademlia.KBuckets[2]))
	assert.Equal(t, 0, len(kademlia.KBuckets[3]))
	assert.Equal(t, 14, len(kademlia.KBuckets[7]))
}

func TestFindClosestKBucket(t *testing.T) {
	idA, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)
	kademlia := NewKademliaTable(idA, "127.0.0.1", "5000")

	lns := prepareTestListedNodes()
	for _, lnI := range lns {
		kademlia.Update(lnI)
	}

	if debug {
		fmt.Println(kademlia)
	}

	idB, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d5")
	assert.Nil(t, err)

	k, err := kademlia.FindClosestKBucket(idB)
	assert.Nil(t, err)
	assert.Equal(t, 2, k)

	idB, err = IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a1000")
	assert.Nil(t, err)

	// the theorical KBucket should be 3
	k = kademlia.KBucket(idB)
	assert.Equal(t, 3, k)

	// while the real KBucket (as the 3 is empty), should be 2
	k, err = kademlia.FindClosestKBucket(idB)
	assert.Nil(t, err)
	assert.Equal(t, 2, k)

	idB, err = IDFromString("0fd85ddddf15aeec2d5d8b01b013dbc000000000")
	assert.Nil(t, err)

	// the theorical KBucket should be 3
	k = kademlia.KBucket(idB)
	assert.Equal(t, 5, k)

	// while the real KBucket (as the 3 is empty), should be 2
	k, err = kademlia.FindClosestKBucket(idB)
	assert.Nil(t, err)
	assert.Equal(t, 7, k)
}
