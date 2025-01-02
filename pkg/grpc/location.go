package grpc

import (
	"context"

	"github.com/goodfoodcesi/tracking-api/pkg/logging"
	pb "github.com/goodfoodcesi/tracking-api/pkg/tracking"
)

func (s *LocationServer) SendLocation(ctx context.Context, location *pb.Location) (*pb.LocationResponse, error) {
	logger := logging.FromContext(ctx)

	logger.Info("Location received",
		"latitude", location.Latitude,
		"longitude", location.Longitude,
		"order_id", location.OrderId,
		"driver_id", location.DriverId,
		"timestamp", location.Timestamp,
	)

	return &pb.LocationResponse{
		Success: true,
		Message: "Location received",
	}, nil
}
