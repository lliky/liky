package main

import (
	"flag"
	"fmt"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	port := flag.Int("port", 0, "the sever port")
	flag.Parse()
	log.Printf("start server on port %d", *port)
	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("image")

	laptopServer := service.NewLaptopServer(laptopStore, imageStore)
	grpcServer := grpc.NewServer()

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)

	address := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
