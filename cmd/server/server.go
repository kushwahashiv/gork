package main

import (
	"sync"

	"context"
	"time"

	"github.com/gork-io/gork/transformers/gateways"
	"github.com/pkg/errors"
)

// NewServer creates a new instance of Server.
func NewServer(options ...ServerOption) (server *Server, err error) {

	server = &Server{}
	for _, option := range options {
		option(server)
	}

	if len(gateways) == 0 {
		return nil, errors.New("no gateways are given")
	}

	return
}

// Server is a container that holds and manages all application gateways.
type Server struct {
	gateways []gateways.Gateway // registered gateways
}

// Start launches registered gateways.
func (srv *Server) Start() (err error) {

	errsChan := make(chan error)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Start all gateways
	for _, gateway := range srv.gateways {
		go func(gateway gateways.Gateway) {
			err := gateway.Start()
			if err != nil {
				errsChan <- err
			}
		}(gateway)
	}

	// Wait for error or context
	select {
	case err = errors.Wrap(<-errsChan, "gateway failed to start"):
		return
	case <-ctx.Done():
		return
	}
}

// Stop asks registered gateways to stop gracefully.
func (srv *Server) Stop() {

	wg := sync.WaitGroup{}

	// Stop all gateways
	wg.Add(len(srv.gateways))
	for _, gateway := range srv.gateways {
		go func(gateway gateways.Gateway) {
			gateway.Stop()
			wg.Done()
		}(gateway)
	}
	wg.Wait()
}

// ServerOption is used to set custom server options.
type ServerOption func(srv *Server)

// ServerWithGateways appends given gateways to the server.
func ServerWithGateways(gateways ...gateways.Gateway) (option ServerOption) {
	return func(srv *Server) {
		srv.gateways = append(srv.gateways, gateways...)
	}
}
