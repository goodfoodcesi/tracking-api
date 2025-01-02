package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/goodfoodcesi/tracking-api/pkg/tracking"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedLocationServiceServer
}

func (s *server) SendLocation(ctx context.Context, in *pb.Location) (*pb.LocationResponse, error) {
	fmt.Println("Location received", in)
	return &pb.LocationResponse{
		Success: true,
		Message: "Location received",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen on port 8080: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterLocationServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
