#!/bin/bash

sudo rm -rf /tmp/mygro 
GOPROXY=http://localhost:8080 GOPATH=/tmp/mygo go get -u github.com/rs/zerolog

