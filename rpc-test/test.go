package main

import (
	"fmt"
	"go-dht/kademlia"
	"log"
	"net/rpc"
)

// Utility to test the Node RPC methods

func prepareTestListedNodes() []kademlia.ListedNode {
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
	var lns []kademlia.ListedNode
	for i := 0; i < len(lnIDs); i++ {
		idI, err := kademlia.IDFromString(lnIDs[i])
		if err != nil {
			panic(err)
		}
		lnI := kademlia.ListedNode{
			ID:   idI,
			Addr: "",
		}
		lns = append(lns, lnI)
	}
	return lns
}

func main() {
	lns := prepareTestListedNodes()

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Fatal("Connection error: ", err)
	}

	var reply kademlia.ListedNode
	for _, ln := range lns {
		err = client.Call("Node.Ping", ln, &reply)
		if err != nil {
			panic(err)
		}
		fmt.Println(reply)
	}
}
