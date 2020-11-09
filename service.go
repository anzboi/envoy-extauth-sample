package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/anzboi/envoy-extauth-sample/pkg/hellopb"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	svr := grpc.NewServer(
		grpc.ChainUnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			md, _ := metadata.FromIncomingContext(ctx)
			log.Println("headers:", md)
			log.Println("method:", info.FullMethod)
			return handler(ctx, req)
		}),
		grpc.ChainStreamInterceptor(func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			md, _ := metadata.FromIncomingContext(ss.Context())
			log.Println("headers:", md)
			log.Println("Stream started:", info.FullMethod)
			return handler(srv, ss)
		}),
	)
	impl := GreeterService{}
	hellopb.RegisterGreeterServer(svr, impl)
	reflection.Register(svr)

	mux := runtime.NewServeMux()
	hellopb.RegisterGreeterHandlerServer(context.Background(), mux, impl)

	log.Println("Starting server")

	errCh := make(chan error)
	go func() { errCh <- svr.Serve(lis) }()
	go func() { errCh <- http.ListenAndServe(":9090", mux) }()
	log.Fatal(<-errCh)
}

type GreeterService struct{}

func (g GreeterService) GetGreeting(ctx context.Context, req *hellopb.GetGreetingRequest) (*hellopb.Greeting, error) {
	return &hellopb.Greeting{
		Message: fmt.Sprint("Hello %s, how is it going?", req.GetName()),
	}, nil
}
