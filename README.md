# go-dht [![Go Report Card](https://goreportcard.com/badge/github.com/arnaucube/go-dht)](https://goreportcard.com/report/github.com/arnaucube/go-dht) [![Build Status](https://travis-ci.org/arnaucube/go-dht.svg?branch=master)](https://travis-ci.org/arnaucube/go-dht)

Kademlia DHT Go implementation.

Following the specification from
- https://pdos.csail.mit.edu/~petar/papers/maymounkov-kademlia-lncs.pdf
- http://xlattice.sourceforge.net/components/protocol/kademlia/specs.html

## Run
To run a node:
```
go run main.go --config config.test0.yaml --debug start
```


## Test
- Scenario:
```
+--+           +--+
|n0+-----------+n1|
+-++           +--+
  |
  |
  |    +--+           +--+
  +----+n2+-----------+n3|
       +--+           +--+
```

- To run 4 test nodes inside a tmux session:
```
bash run-dev-nodes.sh
```

Using the `test.go` in the `rpc-test` directory:

- calls to the node to perform lookups
	- `go run test.go -find`
		- performs an `admin` call to `Find` node, to the `n0`, asking about the `n3`
- calls to simulate kademlia protocol rpc calls
	- `go run test.go -ping`
		- performs the `PING` call
	- `go run test.go -findnode`
		- performs the `FIND_NODE` call
	- `go run test.go -findvalue`
		- performs the `FIND_VALUE` call
	- `go run test.go -store`
		- performs the `STORE` call

