package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "example.com/grpc_assessment/proto"
	"example.com/grpc_assessment/utils"

	"google.golang.org/grpc"
)

// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative ./proto/user.proto

var (
	port = flag.Int("port", 50051, "gRPC server port")
)

type server struct {
	pb.UnimplementedUserServiceServer
}

var User map[string]string

func (*server) Register(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	fmt.Println("Register")
	reqUser := req.GetUser()

	_, ok := User[reqUser.Username]
	if !ok {
		User[reqUser.Username] = reqUser.Password
		return &pb.CreateUserResponse{
			Message: "User has been created successfully",
		}, nil
	}

	return nil, errors.New("user already exists")
}

func (*server) Login(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	fmt.Println("Login")
	reqUser := req.GetUser()
	password, ok := User[reqUser.Username]
	if ok && password == reqUser.Password {
		return &pb.LoginUserResponse{
			Token: utils.GenerateJWT(reqUser.Username, reqUser.Password),
		}, nil
	}
	return nil, errors.New("invalid credentails")
}

func (*server) Logout(ctx context.Context, req *pb.LogoutUserRequest) (*pb.LogoutUserResponse, error) {
	fmt.Println("Logout")
	reqUser := req.GetToken()
	claims, ok := utils.DecodeJwt(reqUser)
	if ok {
		password, ok := User[claims["usernmae"].(string)]
		if ok && password == claims["password"].(string) {
			delete(User, claims["usernmae"].(string))
			return &pb.LogoutUserResponse{
				Message: "User has been logout successfully",
			}, nil
		}
		return nil, errors.New("invalid token")
	}
	return nil, errors.New("invalid token")
}

func main() {
	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))

	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &server{})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve : %v", err)
	}
}
