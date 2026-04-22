package main

import (
	"chickchirick-auth/cmd/config"
	"chickchirick-auth/cmd/factory"
	"chickchirick-auth/cmd/service"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	factory.InitViper()

	appConfig := config.AppConfiguration{}.NewAppConfiguration()

	dbDecorator := service.InitORM(&appConfig.DatabaseConfig)
	defer dbDecorator.CloseDB()

	redisDecorator := service.InitRedis(appConfig.RedisConfig)
	defer redisDecorator.RedisClose()

	//TODO: если не потребуется - удалить
	//httpClient := factory.InitHttpClient()

	errChan := make(chan error, 1)
	grpcServer := factory.BuildAndServeGRPC()

	go func() {
		if err := factory.BuildAndServe(dbDecorator, redisDecorator); err != nil {
			errChan <- fmt.Errorf("http server error: %w", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		fmt.Printf("fatal error: %v\n", err)
	case sig := <-quit:
		fmt.Printf("received signal: %v\n", sig)
	}

	fmt.Println("shutting down...")
	grpcServer.GracefulStop()
	fmt.Println("stopped.")
}
