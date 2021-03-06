package kademlia

import (
	"encoding/hex"
	"encoding/json"
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

func TestIDMarshalers(t *testing.T) {
	id, err := IDFromString("0fd85ddddf15aeec2d5d8b01b013dbca030a18d7")
	assert.Nil(t, err)

	idStr, err := json.Marshal(id)
	assert.Nil(t, err)
	assert.Equal(t, "\"0fd85ddddf15aeec2d5d8b01b013dbca030a18d7\"", string(idStr))

	var idParsed ID
	err = json.Unmarshal(idStr, &idParsed)
	assert.Nil(t, err)
	assert.Equal(t, id, idParsed)

	var idParsed2 ID
	err = json.Unmarshal([]byte("\"0fd85ddddf15aeec2d5d8b01b013dbca030a18d7\""), &idParsed2)
	assert.Nil(t, err)
	assert.Equal(t, id, idParsed2)
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

func TestHashData(t *testing.T) {
	h := HashData([]byte("test data"))
	assert.Equal(t, "916f0027a575074ce72a331777c3478d6513f786", hex.EncodeToString(h[:]))
}
