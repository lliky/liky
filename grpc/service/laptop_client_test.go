package service_test

import (
	"context"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/sample"
	"github.com/liky/grpc/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"log"
	"net"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopServer, serverAddress := startTestLaptopServer(t)
	laptopClient := newTestLaptopClient(t, serverAddress)

	laptop := sample.NewLaptop()
	expectedId := laptop.Id
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}
	log.Printf("serverAddress: %s, laptopClient: %v", serverAddress, laptopClient)
	rsp, err := laptopClient.CreateLaptop(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, rsp)
	require.Equal(t, expectedId, rsp.Id)

	// check that the laptop is saved to the store
	other, err := laptopServer.Store.Find(rsp.Id)
	require.NoError(t, err)
	require.NotNil(t, other)

	// check that saved laptop is the same as the one we send\
	requireSameLaptop(t, laptop, other)
}

func startTestLaptopServer(t *testing.T) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(service.NewInMemoryLaptopStore())

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return laptopServer, listener.Addr().String()
}

func newTestLaptopClient(t *testing.T, serverAddress string) pb.LaptopServiceClient {
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	require.NoError(t, err)
	return pb.NewLaptopServiceClient(conn)
}

func requireSameLaptop(t *testing.T, laptop1, laptop2 *pb.Laptop) {
	//todo
	//require.Equal(t, json1, json2)
}
