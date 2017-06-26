#!/bin/bash

#GOROOT=/usr/local/go
GOPATH=$PWD
go build -o ab -gcflags "-N -l" main