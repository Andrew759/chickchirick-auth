package factory

import (
	internalGrpc "chickchirick-auth/internal/controller/service/grpc"
	authGRPC "chickchirick-auth/internal/gen/auth"
	"fmt"
	"log"
	"net"

	_ "net/http/pprof"

	"google.golang.org/grpc"
)

func BuildAndServeGRPC() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	authGRPC.RegisterAuthServiceServer(s, &internalGrpc.AuthGRPCController{})

	fmt.Println("Auth gRPC server is running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
