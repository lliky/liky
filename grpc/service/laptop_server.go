package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/liky/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// LaptopServer is the server that provides laptop services
type LaptopServer struct {
	pb.UnimplementedLaptopServiceServer
	Store LaptopStore
}

// NewLaptopServer returns a new LaptopServer
func NewLaptopServer(store LaptopStore) *LaptopServer {
	return &LaptopServer{
		Store: store,
	}
}

// CreateLaptop is a unary RPC to create a new laptop
func (server *LaptopServer) CreateLaptop(ctx context.Context, req *pb.CreateLaptopRequest) (*pb.CreateLaptopResponse, error) {
	lapTop := req.GetLaptop()
	log.Printf("receive a create laptop request with id: %s\n", lapTop.Id)
	if len(lapTop.Id) > 0 {
		// check if it's a valid UUID
		_, err := uuid.Parse(lapTop.Id)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "laptop ID is not a valid UUID: %v", err)
		}
	} else {
		id, err := uuid.NewRandom()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "cannot generate a new laptop ID: %v", err)
		}
		lapTop.Id = id.String()
	}
	// save the laptop to store
	err := server.Store.Save(lapTop)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, ErrAlreadyExists) {
			code = codes.AlreadyExists
		}
		return nil, status.Errorf(code, "cannot save laptop to the store: %v", err)
	}

	log.Printf("saved laptop with id: %s", lapTop.Id)

	res := &pb.CreateLaptopResponse{
		Id: lapTop.Id,
	}
	return res, nil
}
