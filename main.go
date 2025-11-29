package main

import (
	converter "curr-conv-api/converter"
	pb "curr-conv-api/proto"
	routes "curr-conv-api/routes"

	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverPort uint16 = 8080
	grpcPort   uint16 = 50051
)

func main() {
	manager := converter.ConvManager{}

	fmt.Println("Initializing manager...")
	if err := manager.Init(); err != nil {
		log.Fatalf("Failed to initialize manager: %v.", err)
		return
	}

	go func() {
		server := routes.CreateGRPCService(&manager)
		err := server.StartService(grpcPort)

		if err != nil {
			log.Fatalf("An error occured when trying to start GRPC Server: %v\n", err)
			return
		}
	}()

	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", grpcPort), grpc.WithTransportCredentials(insecure.NewCredentials())) // client port
	if err != nil {
		log.Fatalf("An error occured when trying to make new grpc client: %v\n", err)
		return
	}

	defer conn.Close()

	client := pb.NewConverterServiceClient(conn)
	httpServer := routes.CreateHTTPServer(client)

	if err := httpServer.Listen(serverPort); err != nil {
		log.Fatalf("An error occured when trying to make new grpc client: %v\n", err)
		return
	}
}
