#!/bin/sh

SESSION='go-dht'

tmux new-session -d -s $SESSION
tmux split-window -d -t 0 -v
tmux split-window -d -t 0 -h
tmux split-window -d -t 2 -h


tmux send-keys -t 0 'go run main.go --config config.test0.yaml --debug start' enter
sleep 2
tmux send-keys -t 1 'go run main.go --config config.test1.yaml --debug start' enter
tmux send-keys -t 2 'go run main.go --config config.test2.yaml --debug start' enter
sleep 1
tmux send-keys -t 3 'go run main.go --config config.test3.yaml --debug start' enter

tmux attach
