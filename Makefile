all: proto

proto:
	protoc \
        --proto_path=${GOPATH}/src/ \
        --proto_path ./transformers/gateways/grpc/proto/protobuf/ \
        --gofast_out=plugins=grpc:./transformers/gateways/grpc/proto/ \
        ./transformers/gateways/grpc/proto/protobuf/*.proto