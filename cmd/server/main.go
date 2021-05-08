package main

import (
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/user_service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	log.Println("server is started")

	server := grpc.NewServer()
	grpcServer := &user_service.GRPCServer{}
	api.RegisterUserServiceServer(server, grpcServer)
	listener, error := net.Listen("tcp", ":8080")
	if error != nil {
		log.Fatalln(error)
	}
	if error := server.Serve(listener); error != nil {
		log.Fatalln(error)
	}

}
