package service_test

import (
	"bufio"
	"context"
	"fmt"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/sample"
	"github.com/liky/grpc/service"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestClientCreateLaptop(t *testing.T) {
	t.Parallel()
	laptopStore := service.NewInMemoryLaptopStore()
	serverAddress := startTestLaptopServer(t, laptopStore, nil, nil)
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
	other, err := laptopStore.Find(rsp.Id)
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
	serverAddress := startTestLaptopServer(t, store, nil, nil)
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

func TestClientUploadImage(t *testing.T) {
	t.Parallel()
	testImageFolder := "../image"

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore(testImageFolder)

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, imageStore, nil)
	laptopClient := newTestLaptopClient(t, serverAddress)

	imagePath := fmt.Sprintf("%s/abc.png", testImageFolder)
	file, err := os.Open(imagePath)
	require.NoError(t, err)
	defer file.Close()

	stream, err := laptopClient.UploadImage(context.Background())
	require.NoError(t, err)

	imageType := filepath.Ext(imagePath)

	req := &pb.UploadImageRequest{
		Data: &pb.UploadImageRequest_Info{
			Info: &pb.ImageInfo{
				LaptopId:  laptop.GetId(),
				ImageType: imageType,
			},
		},
	}
	err = stream.Send(req)
	require.NoError(t, err)

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)
	size := 0

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		size += n

		req = &pb.UploadImageRequest{Data: &pb.UploadImageRequest_ChunkData{
			ChunkData: buffer[:n],
		}}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)
	require.NotZero(t, res.GetId())
	require.EqualValues(t, size, res.GetSize())

	saveImagePath := fmt.Sprintf("%s/%s%s", testImageFolder, res.GetId(), imageType)
	require.FileExists(t, saveImagePath)
	require.NoError(t, os.Remove(saveImagePath))
}

func TestClientRateLaptop(t *testing.T) {
	t.Parallel()

	laptopStore := service.NewInMemoryLaptopStore()
	ratingStore := service.NewInMemoryRatingStore()

	laptop := sample.NewLaptop()
	err := laptopStore.Save(laptop)
	require.NoError(t, err)

	serverAddress := startTestLaptopServer(t, laptopStore, nil, ratingStore)
	laptopClient := newTestLaptopClient(t, serverAddress)

	stream, err := laptopClient.RateLaptop(context.Background())
	require.NoError(t, err)

	scores := []float64{8, 7.5, 10}
	averages := []float64{8, 7.75, 8.5}

	n := len(scores)
	for i := 0; i < n; i++ {
		req := &pb.RateLaptopRequest{
			LaptopId: laptop.GetId(),
			Score:    scores[i],
		}

		err = stream.Send(req)
		require.NoError(t, err)
	}

	err = stream.CloseSend()
	require.NoError(t, err)

	for idx := 0; ; idx++ {
		res, err := stream.Recv()
		if err == io.EOF {
			require.Equal(t, n, idx)
			return
		}
		require.NoError(t, err)
		require.Equal(t, res.GetLaptopId(), laptop.GetId())
		require.Equal(t, res.RatedCount, uint32(idx+1))
		require.Equal(t, res.GetAverageScore(), averages[idx])
	}
}

func startTestLaptopServer(t *testing.T, laptopStore service.LaptopStore, imageStore service.ImageStore, ratingStore service.RatingStore) string {
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	listener, err := net.Listen("tcp", ":0") // random available port
	require.NoError(t, err)

	go grpcServer.Serve(listener)

	return listener.Addr().String()
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
