package service

import (
	"context"
	"github.com/liky/grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

// AuthServer is the server for authentication
type AuthServer struct {
	pb.UnsafeAuthServiceServer
	userStore  UserStore
	jwtManager *JWTManager
}

// Login is a unary RPC to login user
func (server *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("req: %v", req)
	user, err := server.userStore.Find(req.GetUsername())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot find user: %v", err)
	}
	log.Printf("user: %v", user)
	if user == nil || user.IsCorrectPassword(req.GetPassword()) {
		return nil, status.Errorf(codes.Internal, "incorrect username/password")
	}

	token, err := server.jwtManager.Generate(user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot generate access token")
	}
	return &pb.LoginResponse{
		AccessToken: token,
	}, nil
}

// NewAuthServer returns a new auth server
func NewAuthServer(userStore UserStore, jwtManager *JWTManager) *AuthServer {
	return &AuthServer{userStore: userStore, jwtManager: jwtManager}
}
