package control

import (
	"context"
	"fmt"
	"net"

	"github.com/kgretzky/evilginx2/log"
	"github.com/kgretzky/evilginx2/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	controlService *ControlService
	server         *grpc.Server
	port           string
}

func NewGRPCServer(controlService *ControlService, port string) *GRPCServer {
	return &GRPCServer{
		controlService: controlService,
		port:           port,
	}
}

func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", s.port, err)
	}

	s.server = grpc.NewServer()
	proto.RegisterProxyControlServiceServer(s.server, s.controlService)

	log.Info("Starting gRPC control service on port %s", s.port)
	return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
	}
}
