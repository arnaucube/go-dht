package kademlia

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
)

type ID [B]byte

func NewID() (ID, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return ID{}, err
	}
	var id ID
	copy(id[:], b[:B])
	return id, nil
}

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

func IDFromString(s string) (ID, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return ID{}, err
	}
	var id ID
	copy(id[:], b[:B])
	return id, nil
}

func (idA ID) Equal(idB ID) bool {
	return bytes.Equal(idA[:], idB[:])
}

// Cmp returns true if idA > idB
func (idA ID) Cmp(idB ID) bool {
	for i := 0; i < len(idA); i++ {
		if idA[i] != idB[i] {
			return idA[i] > idB[i]
		}
	}
	return false
}

func (idA ID) Distance(idB ID) ID {
	var d ID
	for i := 0; i < B; i++ {
		d[i] = idA[i] ^ idB[i]
	}
	return d
}