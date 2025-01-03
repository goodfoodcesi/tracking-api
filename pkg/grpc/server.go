package grpc

import (
	"github.com/goodfoodcesi/tracking-api/pkg/config"
	"github.com/goodfoodcesi/tracking-api/pkg/logging"
	pb "github.com/goodfoodcesi/tracking-api/pkg/tracking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	grpctrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/google.golang.org/grpc"
)

type LocationServer struct {
	pb.UnimplementedLocationServiceServer
	logManager *logging.LogManager
}

func NewServer(logManager *logging.LogManager, cfg config.Config) *grpc.Server {
	interceptors := []grpc.UnaryServerInterceptor{
		logManager.GrpcInterceptor(),
		JWTInterceptor(cfg),
	}

	// Ajouter le tracing uniquement en prod
	if cfg.Env != "dev" {
		interceptors = append(interceptors, grpctrace.UnaryServerInterceptor())
	}

	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(interceptors...),
	)

	locationServer := &LocationServer{
		logManager: logManager,
	}

	pb.RegisterLocationServiceServer(server, locationServer)
	reflection.Register(server)
	return server
}
