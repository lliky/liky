package service_test

import (
	"context"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/sample"
	"github.com/liky/grpc/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopServer, serverAddress := startTestLaptopServer(t, service.NewInMemoryLaptopStore())
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

func TestClientSearchLaptop(t *testing.T) {
	t.Parallel()
	filter := &pb.Filter{
		MaxPriceUsd: 2000,
		MinCpuCores: 4,
		MinCpuGhz:   2.2,
		MinRam: &pb.Memory{
			Value: 8,
			Unit:  pb.Memory_GIGABYTE,
		}}

	store := service.NewInMemoryLaptopStore()
	expectedIDs := make(map[string]bool)
	for i := 0; i < 6; i++ {
		laptop := sample.NewLaptop()
		switch i {
		case 0:
			laptop.PriceUsd = 2500
		case 1:
			laptop.Cpu.NumberCores = 2
		case 2:
			laptop.Cpu.MinGhz = 2.0
		case 3:
			laptop.Ram = &pb.Memory{Value: 30, Unit: pb.Memory_MEGABYTE}
		case 4:
			laptop.PriceUsd = 1900
			laptop.Cpu.NumberCores = 4
			laptop.Cpu.MinGhz = 2.6
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{Value: 16, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		case 5:
			laptop.PriceUsd = 2000
			laptop.Cpu.NumberCores = 6
			laptop.Cpu.MinGhz = 2.99
			laptop.Cpu.MaxGhz = 4.5
			laptop.Ram = &pb.Memory{Value: 32, Unit: pb.Memory_GIGABYTE}
			expectedIDs[laptop.Id] = true
		}
		err := store.Save(laptop)
		require.NoError(t, err)
	}
	_, serverAddress := startTestLaptopServer(t, store)
	laptopClient := newTestLaptopClient(t, serverAddress)

	req := &pb.SearchLaptopRequest{Filter: filter}
	stream, err := laptopClient.SearchLaptop(context.Background(), req)
	require.NoError(t, err)

	found := 0
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Contains(t, expectedIDs, res.Laptop.GetId())
		found += 1
	}

	require.Equal(t, len(expectedIDs), found)
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore) (*service.LaptopServer, string) {
	laptopServer := service.NewLaptopServer(laptopStore)

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
