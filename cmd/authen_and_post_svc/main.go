package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/leducanh112/newsfeed/configs"
	"github.com/leducanh112/newsfeed/internal/app/authen_and_post_svc"
	"github.com/leducanh112/newsfeed/pkg/types/proto/pb/authen_and_post"
)

// command line flag
var (
	path = flag.String("c", "config.yml", "config path for this service")
)

func main() {
	flag.Parse()

	// Start authenticate and post service
	conf, err := configs.GetAuthenticateAndPostConfig(*path)
	if err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	log.Printf("starts with config: %+v\n", conf)

	service, err := authen_and_post_svc.NewAuthenticateAndPostService(conf)
	if err != nil {
		log.Fatalf("failed to init server %s", err)
	}

	addr := fmt.Sprintf("0.0.0.0:%d", conf.Port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	authen_and_post.RegisterAuthenticateAndPostServer(grpcServer, service)

	log.Printf("grpc server starting on %+v\n", addr)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("server stopped %v", err)
	}
}
