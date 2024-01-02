package main

import (
	"context"
	"flag"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server :%s", *serverAddress)
	conn, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	laptopClient := pb.NewLaptopServiceClient(conn)

	laptop := sample.NewLaptop()
	req := &pb.CreateLaptopRequest{
		Laptop: laptop,
	}

	rsp, err := laptopClient.CreateLaptop(context.Background(), req)
	if err != nil {
		s, ok := status.FromError(err)
		if ok && s.Code() == codes.AlreadyExists {
			log.Printf("laptop already exists")
		} else {
			log.Fatal("cannot create laptop")
		}
		return
	}
	log.Printf("created laptop with id: %s", rsp.Id)
}
