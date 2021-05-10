package main

import (
	"github.com/fev0ks/UserServiceSC/cmd/server/db"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/service"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	serverPort = ":8080"
	network    = "tcp"
)

func main() {
	log.Println("Starting...")
	db.InitDataBase()
	startServer()
}

func startServer() {
	log.Println("server is started")
	server := grpc.NewServer()
	grpcServer := &service.GRPCServer{}
	api.RegisterUserServiceServer(server, grpcServer)
	listener, err := net.Listen(network, serverPort)
	if err != nil {
		log.Fatalln(err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
