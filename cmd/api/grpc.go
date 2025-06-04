package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/sangketkit01/7-coding-test/internal/db"
	"github.com/sangketkit01/7-coding-test/internal/util"
	"github.com/sangketkit01/7-coding-test/pb"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GRPCService struct {
	pb.UnimplementedSevenCodingTestServer
	model db.MongoClient
}

func (service *GRPCService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	if strings.TrimSpace(req.GetName()) == "" || strings.TrimSpace(req.GetEmail()) == "" || len(strings.TrimSpace(req.GetPassword())) < 8 {
		return nil, errors.New("Name and Email is required and password's length must be atleast 8 characters long.")
	}

	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, errors.New("failed to hash password.")
	}

	newUser := db.User{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: hashedPassword,
	}

	err = service.model.Insert(newUser)
	if err != nil {
		return nil, err
	}

	response := &pb.CreateUserResponse{
		User: &pb.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
		},
	}

	return response, nil

}

func (service *GRPCService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if strings.TrimSpace(req.GetXId()) == "" {
		return nil, errors.New("id is not provided.")
	}

	user, err := service.model.FetchUserByID(req.GetXId())
	if err != nil {
		return nil, err
	}

	response := &pb.GetUserResponse{
		User: &pb.User{
			XId:       user.ID.Hex(),
			Name:      user.Name,
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(user.CreatedAt),
		},
	}

	return response, nil
}

func (app *App) gRPCListen() {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Fatalln("Failed to listen grpc:", err)
	}

	server := grpc.NewServer()

	pb.RegisterSevenCodingTestServer(server, &GRPCService{model: app.model})

	log.Printf("gRPC server started at port: %s\n", gRpcPort)

	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to listen grpc: %v\n", err)
	}
}
