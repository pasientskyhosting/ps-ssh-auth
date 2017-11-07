#!/bin/sh
export GOPATH=$(PWD)/packages/
go get
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ssh-auth.go
