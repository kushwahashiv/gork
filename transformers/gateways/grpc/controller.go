package grpc

import "google.golang.org/grpc"

// Controller is an interface that should be implemented by all GRPC controllers.
type Controller interface {
	// Register registers this controller as a GRPC service implementation.
	Register(server *grpc.Server)
}
