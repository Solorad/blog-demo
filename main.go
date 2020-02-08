package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	gw "github.com/solorad/blog-demo/server/pkg/api/v1" // Update
	"github.com/solorad/blog-demo/server/pkg/blog"
	"github.com/solorad/blog-demo/server/pkg/log"
	"github.com/solorad/blog-demo/server/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "0.0.0.0:50051", "gRPC server endpoint")
)

func main() {
	flag.Parse()
	defer glog.Flush()

	startGRPCServer()
	if err := runReverseProxy(); err != nil {
		glog.Fatal(err)
	}
}

func runReverseProxy() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := gw.RegisterBlogServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)


	log.Infof("GRPC server started on %v", *grpcServerEndpoint)
	if err != nil {
		return err
	}
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Infof("REST server is starting on port 8081")
	return http.ListenAndServe(":8081", mux)
}

func startGRPCServer() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	s := grpc.NewServer(opts...)
	mongoClient := storage.NewMongoClient()
	blogServer := blog.NewBlogServiceServer(mongoClient)
	gw.RegisterBlogServiceServer(s, blogServer)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
}
