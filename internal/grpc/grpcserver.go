package grpc

import (
	"context"
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"microservice-exercise/internal/data_model"
	"microservice-exercise/internal/database"
	"microservice-exercise/internal/transport"

	"google.golang.org/grpc"
)

// PortsProtocolServer is a server handling ports protocol, it is made by the db to send the data to,
// the grpc Server and a listener over the port (tcp)
type PortsProtocolServer struct {
	transport.UnimplementedPortServiceServer
	server   *grpc.Server
	db       database.DatabaseService
	listener net.Listener
}

// NewPortsProtocolServer returns new portsPortocolServer
func NewPortsProtocolServer(serverPort int, portsDb database.DatabaseService, cancel context.CancelFunc) *PortsProtocolServer {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", serverPort))
	if err != nil {
		log.Error(err)
		cancel()
	}

	newServer := &PortsProtocolServer{
		server:   grpcServer,
		db:       portsDb,
		listener: lis,
	}
	transport.RegisterPortServiceServer(grpcServer, newServer)
	return newServer
}

// CreateOrUpdatePort method creates or updates port information in database
func (s *PortsProtocolServer) UpdatePort(ctx context.Context, request *data_model.UpdatePortRequest) (*data_model.UpdatePortResponse, error) {
	log.Debug("Got UpdatePort request for: ", request.Key, " port")
	err := s.db.UpdatePort(request.Key, *request.Port)
	return &data_model.UpdatePortResponse{}, err
}

// GracefulStop stops grpc server and listener as well
func (s *PortsProtocolServer) GracefulStop() error {
	log.Info("Please wait while stopping gRPC server...")
	s.server.GracefulStop()
	log.Info("Stopping gRPC done!")
	log.Info("Please wait while stopping server listener...")
	return s.listener.Close()
	log.Info("Stopping server listener done!")
}

// Serve serves the grpc server
func (s *PortsProtocolServer) Serve() error {
	return s.server.Serve(s.listener)
}
