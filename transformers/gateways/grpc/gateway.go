package grpc

import (
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"google.golang.org/grpc"
)

// NewGateway creates a new instance of Gateway.
func NewGateway() (gateway *Gateway) {
	return &Gateway{
		server: grpc.NewServer(
			grpc.UnaryInterceptor(
				grpc_middleware.ChainUnaryServer(
					grpc_recovery.UnaryServerInterceptor(),
				),
			),
			grpc.StreamInterceptor(
				grpc_middleware.ChainStreamServer(
					grpc_recovery.StreamServerInterceptor(),
				),
			),
		),
	}
}

// Gateway is an GRPC implementation of the Gorp gateway.
type Gateway struct {
	server      *grpc.Server // GRPC server instance
	listener    net.Listener // listener to bind to
	controllers []Controller // controllers (GRPC services) to expose
}

func (gtw *Gateway) Name() (name string) {
	return "GRPC"
}

func (gtw *Gateway) Start() (err error) {

	// Register controllers
	for _, controller := range gtw.controllers {
		controller.Register(gtw.server)
	}

	return gtw.server.Serve(gtw.listener)
}

func (gtw *Gateway) Stop() {
	gtw.server.GracefulStop()
}
