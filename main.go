package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	proto "securitymedic/proto"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	port = flag.Int("port", 0, "port for peer")
	name = flag.String("name", "", "name of peer")
)
var hospitalResponse = 0
var sentChunks = 0

type Peer struct {
	proto.UnimplementedSecretServiceServer
	id      int32
	name    string
	clients map[int32]proto.SecretServiceClient
	chunks  map[int32]int
	ctx     context.Context
}

func main() {
	flag.Parse()

	ownPort := int32(*port + 5000)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &Peer{
		id:      ownPort,
		name:    *name,
		clients: make(map[int32]proto.SecretServiceClient),
		chunks:  make(map[int32]int),
		ctx:     ctx,
	}

	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", ownPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serverCert, err := credentials.NewServerTLSFromFile("certificate/server.crt", "certificate/priv.key")
	if err != nil {
		log.Fatalln("failed to create cert", err)
	}


	grpcServer := grpc.NewServer(grpc.Creds(serverCert))
	proto.RegisterSecretServiceServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	for i := 0; i <= 3; i++ {

		port := 5000 + int32(i)

		if port == int32(p.id) {
			continue
		}

		clientCert, err := credentials.NewClientTLSFromFile("certificate/server.crt", "")
		if err != nil {
			log.Fatalln("failed to create cert", err)
		}


		fmt.Printf("Trying to dial: %v\n", port)

		conn, err := grpc.Dial(fmt.Sprintf("localhost:%v", port), grpc.WithTransportCredentials(clientCert), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %s", err)
		}
		defer conn.Close()
		c := proto.NewSecretServiceClient(conn)
		p.clients[port] = c
		fmt.Printf("Dialing: %v\n", port)

	}

	if ownPort != 5000 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Enter number: ")
		for scanner.Scan() {
			input, _ := strconv.ParseInt(scanner.Text(), 10, 32)
			fmt.Printf("Sending: %v\n", input)

			p.ShareChunkToPeers(int32(input))
			break
		}

		// Wait for all chunks to be sent

		for sentChunks < 2 {
			time.Sleep(2 * time.Second)
		}

		
		
		sumChunks := p.SumOfChunks()
		fmt.Printf("%v: Sending my sum to the hospital\n", p.name)

		hospital := p.clients[5000]
		req := &proto.Message{Id: p.id, Content: int32(sumChunks)}
		reply, err := hospital.SendResult(p.ctx, req)

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		fmt.Printf("Response from hospital: %v\n", reply.Response)
		time.Sleep(6 * time.Second)

	} else {
		fmt.Printf("I am the hospital\n")
		for hospitalResponse < 3 {
			time.Sleep(2 * time.Second)
		}
		sum := int(0)
		for _, chunk := range p.chunks {
			sum += chunk
		}
		fmt.Printf("Hospital: The total sum for all the patients is %v\n", sum)


	}

	for {

	}

}
func (p *Peer) SendChonker(ctx context.Context, in *proto.Chunk) (*proto.Response, error) {
	id := in.Id
	p.chunks[id] = int(in.Chunk)
	/* name := ""
	switch id {
	case 5001:
		name = "Alice"
	case 5002:
		name = "Bob"
	case 5003:
		name = "Charlie"
	} */
	sentChunks++
	rString := fmt.Sprintf("Sending %v, to %v", p.chunks[id], p.name)
	return &proto.Response{Response: rString}, nil

}

func (p *Peer) SendResult(ctx context.Context, in *proto.Message) (*proto.Response, error) {
	id := in.Id
	name := ""
	switch id {
	case 5001:
		name = "Alice"
	case 5002:
		name = "Bob"
	case 5003:
		name = "Charlie"
	}
	hospitalResponse++
	p.chunks[id] += int(in.Content)
	rString := fmt.Sprintf("Hospital: I got %v from %v", in.Content, name)
	return &proto.Response{Response: rString}, nil
}

func (p *Peer) ShareChunkToPeers(chunk int32) {

	s1, s2, s3 := SplitMessage(int(chunk))
	p.chunks[p.id] = s1
	fmt.Printf("I am %v, My share is %v\n", p.name, p.chunks[p.id])

	counter := 0
	for id, client := range p.clients {
		if id == 5000 {
			continue
		}
		var chunk int32
		if counter == 0 {
			chunk = int32(s2)
		} else {
			chunk = int32(s3)
		}

		req := &proto.Chunk{Chunk: int32(chunk), Id: p.id}

		resp, err := client.SendChonker(p.ctx, req)
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		fmt.Printf("Response: %v\n", resp)
		counter++
	}

}

func (p *Peer) SumOfChunks() int {
	sum := int(0)
	for _, chunk := range p.chunks {
		sum += chunk
	}
	return sum
}

func SplitMessage(n int) (s1, s2, s3 int) {

	s1 = rand.Intn(n)
	s2 = rand.Intn(n - s1)
	s3 = n - s1 - s2

	return s1, s2, s3
}
