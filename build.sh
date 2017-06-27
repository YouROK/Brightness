#!/bin/bash

GOROOT=/usr/local/go
GOPATH=$PWD
echo "GOPATH=$GOPATH"
go build -o ab -gcflags "-N -l" ./src/main