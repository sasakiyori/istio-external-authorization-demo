#!/bin/bash

# docker hub user name
USERNAME="sasakiyori"

# build binary file
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ext_authz_server authorization.go

# build docker image
docker build -t ext_authz_server:0.0.1 ./

# create docker tag
docker tag ext_authz_server:0.0.1 ${USERNAME}/ext_authz_server:0.0.1

# login
docker login

# push docker image
docker push ${USERNAME}/ext_authz_server:0.0.1

# clear binary file
rm ext_authz_server
