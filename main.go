package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	pb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

var (
	port = flag.Int("port", 50051, "the port to serve on")

	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
)

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) ServerStreamingEcho(req *pb.EchoRequest, stream pb.Echo_ServerStreamingEchoServer) error {
	_ = stream.Send(&pb.EchoResponse{Message: "abc"})

	return nil
}

func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return errMissingMetadata
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.StreamInterceptor(streamInterceptor))

	// Register EchoServer on the server.
	pb.RegisterEchoServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
