package client

import (
	"context"
	"github.com/liky/grpc/pb"
	"google.golang.org/grpc"
	"time"
)

// AuthClient is a client to call authentication RPC
type AuthClient struct {
	service  pb.AuthServiceClient
	username string
	password string
}

// / NewAuthClient returns a new auth client
func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	client := pb.NewAuthServiceClient(cc)
	return &AuthClient{
		service:  client,
		username: username,
		password: password,
	}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pb.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	rsp, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}
	return rsp.GetAccessToken(), nil
}
