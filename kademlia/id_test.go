package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewID(t *testing.T) {
	// get some IDs
	// for i := 0; i < 10; i++ {
	//         id, err := NewID()
	//         assert.Nil(t, err)
	//         fmt.Println(id)
	// }

	idA, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)
	assert.Equal(t, "0fd85ddddf15aeec2d5d8b01b013dbca030a18d7", idA.String())
}

func TestIDCmp(t *testing.T) {
	idA, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)
	idB, err := IDFromString("c48d8b53dbefb609ed4e94d386dd5b22efcb2c5b")
	assert.Nil(t, err)
	assert.True(t, !idA.Cmp(idB))
}

func TestIDDistance(t *testing.T) {
	idA, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)
	idB, err := IDFromString("c48d8b53dbefb609ed4e94d386dd5b22efcb2c5b")
	assert.Nil(t, err)
	assert.Equal(t, "cb55d68e04fa18e5c0131fd236ce80e8ecc1348c", idA.Distance(idB).String())
}
