# !/bin/bash
proto=$1
protoc -I. --plugin=$GOPATH/bin/protoc-gen-go $proto --go_out=.
