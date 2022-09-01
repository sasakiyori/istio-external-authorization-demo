#!/bin/bash

# docker hub user name
USERNAME="sasakiyori"

# build binary file
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ext-authz-server authorization.go

# build docker image
docker build -t ext-authz-server:0.0.1 ./

# create docker tag
docker tag ext-authz-server:0.0.1 ${USERNAME}/ext-authz-server:0.0.1

# login
docker login

# push docker image
docker push ${USERNAME}/ext-authz-server:0.0.1

# clear binary file
rm ext-authz-server
