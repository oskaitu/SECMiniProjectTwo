package main

import (
	"context"
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"

	proto "securitymedic/proto"
)

type Server struct {
	proto.UnimplementedHospitalServiceServer
	name string
	port int
}

var port = flag.Int("port", 0, "server port number")

func main() {

	flag.Parse()

	server := &Server{
		name: "Hospital",
		port: *port,
	}

	go startServer(server)

	for {

	}
}

func startServer(server *Server) {
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("could not create the server %v", err)
	}

	log.Printf("Started server at port: %d/n", server.port)

	proto.RegisterHospitalServiceServer(grpcServer, server)

	serveError := grpcServer.Serve(listener)

	if serveError != nil {
		log.Fatalf("could not serve listener")
	}
}

func (c *Server) SendChunk(ctx context.Context, in *proto.Chunk) (*proto.Response, error) {

	log.Printf("Received: %v", in.Info)
	return &proto.Response{
		Message: "Hello",
		}, nil

}
