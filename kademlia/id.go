package kademlia

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
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

func (id ID) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(id[:])), nil
}

func (id *ID) UnmarshalText(data []byte) error {
	var err error
	var idFromStr ID
	idFromStr, err = IDFromString(string(data))
	if err != nil {
		return err
	}
	copy(id[:], idFromStr[:])
	return nil
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

func HashData(d []byte) ID {
	h := sha256.Sum256(d)
	var r ID
	copy(r[:], h[:20])
	return r
}
