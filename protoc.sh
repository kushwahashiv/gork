#!/usr/bin/env bash

protoc \
    --proto_path ./transformers/gateways/grpc/proto/ \
    --proto_path=./vendor/github.com/mwitkow/go-proto-validators/ \
    --go_out=plugins=grpc:./transformers/gateways/grpc/proto/ \
    --govalidators_out=./transformers/gateways/grpc/proto/ \
    ./transformers/gateways/grpc/proto/*.proto