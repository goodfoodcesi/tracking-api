package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewMultiplexHandler(ginHandler *gin.Engine, grpcServer *grpc.Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && r.Header.Get("content-type") == "application/grpc" {
			grpcServer.ServeHTTP(w, r)
			return
		}
		ginHandler.ServeHTTP(w, r)
	})
}
