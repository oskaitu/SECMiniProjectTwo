package main

import (
	"flag"
	"log"
	"net"
	"strconv"

	"google.golang.org/grpc"
	glog "google.golang.org/grpc/grpclog"

	proto "securitymedic/proto"
)

var grpcLog glog.LoggerV2

/*
type Connection struct {
	stream  proto.SecurityMedic_JoinServer
	user_id string
	name    string
	error   chan error
}
*/


type Server struct {
	proto.UnimplementedSecurityMedicServer
	port int
}

var port = flag.Int("port", 0, "server port number")

func main() {
	flag.Parse()
	
	server := &Server{
		port: *port,
	}

	go serverStart(server)

	for{

	}
}

func serverStart (server *Server){
	grpcServer := grpc.NewServer()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))
	if err != nil {
		log.Fatalf("Error creating the server %v", err)
	}

	grpcLog.Info("Starting server at port :%d",server.port)

	proto.RegisterSecurityMedicServer(grpcServer, server)
	serverError := grpcServer.Serve(listener)

	if serverError != nil {
		log.Fatalf("cold not serve")
	}
}