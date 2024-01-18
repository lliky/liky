package main

import (
	"flag"
	"fmt"
	"github.com/liky/grpc/pb"
	"github.com/liky/grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"time"
)

func seedUser(userStore service.UserStore) error {
	err := createUser(userStore, "admin1", "secret", "admin")
	if err != nil {
		return err
	}
	return createUser(userStore, "user1", "secret", "user")

}

func createUser(userStore service.UserStore, username, password, role string) error {
	user, err := service.NewUser(username, password, role)
	if err != nil {
		return err
	}
	return userStore.Save(user)
}

func accessibleRoles() map[string][]string {
	const path = "/LaptopService/"
	return map[string][]string{
		path + "CreateLaptop": {"admin"},
		path + "UploadImage":  {"admin"},
		path + "RateLaptop":   {"admin", "user"},
	}
}

const (
	secretKey     = "secret"
	tokenDuration = 15 * time.Minute
)

func main() {
	port := flag.Int("port", 0, "the sever port")
	flag.Parse()
	log.Printf("start server on port %d", *port)
	userStore := service.NewInMemoryUserStore()

	err := seedUser(userStore)
	if err != nil {
		log.Fatal("cannot seed users")
	}

	jwtManager := service.NewJWTManager(secretKey, tokenDuration)
	authServer := service.NewAuthServer(userStore, jwtManager)

	laptopStore := service.NewInMemoryLaptopStore()
	imageStore := service.NewDiskImageStore("image")
	ratingStore := service.NewInMemoryRatingStore()
	laptopServer := service.NewLaptopServer(laptopStore, imageStore, ratingStore)

	interceptor := service.NewAuthInterceptor(jwtManager, accessibleRoles())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary()),
		grpc.StreamInterceptor(interceptor.Stream()),
	)

	pb.RegisterLaptopServiceServer(grpcServer, laptopServer)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	reflection.Register(grpcServer)

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
