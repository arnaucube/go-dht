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

To run 3 test nodes inside a tmux session:
```
bash run-dev-nodes.sh
```
