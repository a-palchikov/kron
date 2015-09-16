# !/bin/bash

protoc -I. --plugin=/home/dosler/dev/go/bin/protoc-gen-go servicepb/scheduler.proto servicepb/feedback.proto --go_out=plugins=grpc:.

protoc -I. --plugin=/home/dosler/dev/go/bin/protoc-gen-go storepb/store.proto --go_out=plugins=grpc,Mservicepb/scheduler.proto=github.com/a-palchikov/kron/proto/servicepb:.
