#!/usr/bin/env bash

protoc \
    --proto_path=${GOPATH}/src/ \
    --proto_path ./transformers/gateways/grpc/proto/protobuf/ \
    --proto_path=./vendor/github.com/mwitkow/go-proto-validators/ \
    --gofast_out=plugins=grpc:./transformers/gateways/grpc/proto/ \
    --govalidators_out=gogoimport=true:./transformers/gateways/grpc/proto/ \
    ./transformers/gateways/grpc/proto/protobuf/*.proto