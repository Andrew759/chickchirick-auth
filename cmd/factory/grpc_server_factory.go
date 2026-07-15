package factory

import (
	internalGrpc "chickchirick-auth/internal/controller/service/grpc"
	authGRPC "chickchirick-auth/internal/gen/auth"
	"chickchirick-auth/pkg/chirick_config"
	"errors"
	"fmt"
	"log"
	"net"

	_ "net/http/pprof"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func BuildAndServeGRPC() *grpc.Server {
	lis, err := net.Listen("tcp", viper.GetString(chirick_config.ProtobufServer))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authGRPC.RegisterAuthServiceServer(grpcServer, &internalGrpc.AuthGRPCController{})

	go func() {
		fmt.Printf("Auth gRPC server is running on %s\n", viper.GetString(chirick_config.ProtobufServer))
		if err := grpcServer.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	return grpcServer
}
