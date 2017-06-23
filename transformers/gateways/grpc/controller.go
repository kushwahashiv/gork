package grpc

import "google.golang.org/grpc"

// Controller
type Controller interface {
	Register(server *grpc.Server)
}
