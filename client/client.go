package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"net"
	"os"
	"strconv"

	proto "securitymedic/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var amount_of_clients = 3

var (
	cid        = flag.Int("id", 0, "client ID")
	serverPort = flag.Int("sport", 0, "server port")
	clientName = flag.String("clientname", "", "client name")
)

type Client struct {
	proto.UnimplementedSecretServiceServer
	id         int
	name       string
	clientPort int
	clients    map[int32]proto.SecretServiceClient
	ctx        context.Context
}

func main() {

	flag.Parse()

	ownPort := int(*cid) + 5001

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := &Client{
		id:         *cid,
		name:       *clientName,
		clientPort: ownPort,
		ctx:        ctx,
		clients:    make(map[int32]proto.SecretServiceClient),
	}

	
	PeerConnection(client)
	ServerConnection(client)


	for {

	}

}

func PeerConnection(client *Client) {
	

	list, err := net.Listen("tcp", ":"+strconv.Itoa(client.clientPort))
	if err != nil {
		log.Fatalf("could not listen to port %d", client.clientPort)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterSecretServiceServer(grpcServer, client)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("could not serve %v", err)
		}
	}()

	for i := 0; i < amount_of_clients; i++ {
		port := int32(5001) + int32(i)

		if port == int32(client.clientPort) {
			continue
		}

		var conn *grpc.ClientConn
		conn, err := grpc.Dial("localhost:"+strconv.Itoa(int(port)), grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect to port %d", port)
		}

		defer conn.Close()
		client.clients[port] = proto.NewSecretServiceClient(conn)
	}

	

	
}

func ServerConnection(client *Client) {
	serverConnection, _ := connectToServer()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		log.Printf("Client %d: %s", client.id, input)

		chunkReturnMessage, err := serverConnection.SendChunk(context.Background(),
			&proto.Chunk{
				Info: input,
			})
	
		if err != nil {
			log.Printf("Error sending chunk to server: %v", err.Error())
		} else {
			log.Printf("Server response: %s", chunkReturnMessage.Message)
		}
	}

}

func connectToServer() (proto.HospitalServiceClient, error) {
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("connected to the server at port %d \n", *serverPort)
	}
	return proto.NewHospitalServiceClient(conn), nil
}




