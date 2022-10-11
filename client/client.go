package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/adamjhr/time/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var clientName = flag.String("name", "default", "Senders name")
var serverPort = flag.String("server", "5400", "Server port")
var server pb.TimeClient
var ServerConnection *grpc.ClientConn

func main() {
	flag.Parse()

	ConnectToServer()
	defer ServerConnection.Close()

	for {
		time.Sleep(5 * time.Second)
		fmt.Println(getTime())
	}

}

func getTime() string {

	request := &pb.TimeRequest{}

	ack, err := server.GetCurrentTime(context.Background(), request)

	if err != nil {
		log.Printf("Client %s: no response from the server, attempting to reconnect", *clientName)
		log.Println(err)
	}

	return ack.Time

}

func ConnectToServer() {

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	log.Printf("client %s: Attempts to dial on port %s\n", *clientName, *serverPort)
	conn, err := grpc.Dial(fmt.Sprintf(":%s", *serverPort), opts...)
	if err != nil {
		log.Printf("Fail to Dial : %v", err)
		return
	}

	server = pb.NewTimeClient(conn)
	ServerConnection = conn
}
