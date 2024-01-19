package main

import (
	"flag"
	"fmt"
	"github.com/liky/grpc/client"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/sample"
	"google.golang.org/grpc"
	"log"
	"strings"
	"time"
)

func testCreateLaptop(laptopClient *client.LaptopClient) {
	laptopClient.CreateLaptop(sample.NewLaptop())
}

func testSearchLaptop(laptopClient *client.LaptopClient) {
	for i := 0; i < 10; i++ {
		laptopClient.CreateLaptop(sample.NewLaptop())
	}

	filter := &pb.Filter{
		MaxPriceUsd: 3000,
		MinCpuCores: 4,
		MinCpuGhz:   2.5,
		MinRam: &pb.Memory{
			Value: 8, Unit: pb.Memory_GIGABYTE,
		}}

	laptopClient.SearchLaptop(filter)
}

func testUploadImage(laptopClient *client.LaptopClient) {
	laptop := sample.NewLaptop()
	laptopClient.CreateLaptop(laptop)
	laptopClient.UploadImage(laptop.Id, "image/abc.png")
}

func testRateLaptop(laptopClient *client.LaptopClient) {
	n := 3
	laptopIDs := make([]string, n)
	for i := 0; i < 3; i++ {
		laptop := sample.NewLaptop()
		laptopIDs[i] = laptop.GetId()
		laptopClient.CreateLaptop(laptop)
	}

	scores := make([]float64, n)
	for {
		fmt.Print("rate laptop(y/n)")
		var answer string
		fmt.Scan(&answer)
		if strings.ToLower(answer) != "y" {
			break
		}
		for i := 0; i < n; i++ {
			scores[i] = sample.RandomLaptopScore()
		}
		err := laptopClient.RateLaptop(laptopIDs, scores)
		if err != nil {
			log.Fatal(err)
		}
	}
}

const (
	username        = "admin1"
	password        = "secret"
	refreshDuration = 30 * time.Second
)

func authMethods() map[string]bool {
	const path = "/LaptopService/"
	return map[string]bool{
		path + "CreateLaptop": true,
		path + "UploadImage":  true,
		path + "RateLaptop":   true,
	}
}
func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server :%s", *serverAddress)
	cc1, err := grpc.Dial(*serverAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := client.NewAuthClient(cc1, username, password)
	interceptor, err := client.NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.Dial(*serverAddress, grpc.WithInsecure(), grpc.WithUnaryInterceptor(interceptor.Unary()), grpc.WithStreamInterceptor(interceptor.Stream()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}
	laptopClient := client.NewLaptopClient(cc2)
	testRateLaptop(laptopClient)
}
