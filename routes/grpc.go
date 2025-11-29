package routes

import (
	"context"
	converter "curr-conv-api/converter"
	pb "curr-conv-api/proto"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedConverterServiceServer
	Manager *converter.ConvManager
}

func CreateGRPCService(manager *converter.ConvManager) *GRPCServer {
	return &GRPCServer{Manager: manager}
}

func (server *GRPCServer) GetCurrencyList(ctx context.Context, req *pb.GetCurrencyListRequest) (*pb.GetCurrencyListResponse, error) {
	currencyList := server.Manager.GetCurrencies()
	return &pb.GetCurrencyListResponse{List: currencyList}, nil
}

func (server *GRPCServer) DoConvert(ctx context.Context, req *pb.DoConvertRequest) (*pb.DoConvertResponse, error) {
	from := req.GetFrom()
	to := req.GetTo()
	amount := req.GetAmount()

	//fmt.Println(from, to, amount)

	result, err := server.Manager.Convert(from, to, amount)
	if err != nil {
		return nil, err
	}

	return &pb.DoConvertResponse{Result: result}, nil
}

func (server *GRPCServer) StartService(grpcPort uint16) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Failed to start server in port %d: %v\n", grpcPort, err)
		return err
	}

	gRpcServer := grpc.NewServer()
	todoService := CreateGRPCService(server.Manager)

	pb.RegisterConverterServiceServer(gRpcServer, todoService)
	fmt.Printf("Service is running in port %d.\n", grpcPort)

	if err := gRpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to run server in port %d: %v\n", grpcPort, err)
		return err
	}

	return nil
}
