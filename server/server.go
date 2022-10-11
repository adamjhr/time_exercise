package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/adamjhr/time/proto"

	"google.golang.org/grpc"
)

type TimeServer struct {
	pb.UnimplementedTimeServer
	name  string
	port  string
	mutex sync.Mutex
}

var serverName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("port", "5400", "Senders port")

func main() {
	flag.Parse()

	go launchServer()

	for {
		time.Sleep(time.Second * 1)
	}
}

func launchServer() {
	list, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", *serverPort))

	if err != nil {
		log.Printf("Server %s: Failed to listen on port %s: %v", *serverName, *serverPort, err)
		return
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	server := &TimeServer{
		name: *serverName,
		port: *serverPort,
	}

	pb.RegisterTimeServer(grpcServer, server)

	log.Printf("Server %s: Listening at %v\n", *serverName, list.Addr())

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}

func (s *TimeServer) GetTime(ctx context.Context, request *pb.TimeRequest) (*pb.TimeReply, error) {

	time := time.Now().String()

	return &pb.TimeReply{Time: time}, nil

}
