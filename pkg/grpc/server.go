package grpc

import (
	pb "github.com/goodfoodcesi/tracking-api/pkg/tracking"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type LocationServer struct {
	pb.UnimplementedLocationServiceServer
}

func NewServer() *grpc.Server {
	server := grpc.NewServer()
	pb.RegisterLocationServiceServer(server, &LocationServer{})
	reflection.Register(server)
	return server
}
