#!/bin/bash
cd ../cmd
go get github.com/go-playground/webhooks/v6
go get github.com/google/go-github/v45
go get golang.org/x/oauth2

CGO_ENABLED=0 go build -o $1 .