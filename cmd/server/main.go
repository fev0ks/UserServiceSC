package main

import (
	"database/sql"
	"fmt"
	api "github.com/fev0ks/UserServiceSC/pkg/api"
	"github.com/fev0ks/UserServiceSC/pkg/user_service"
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	host          = "localhost"
	serverPort    = ":8080"
	network       = "tcp"
	migrationsDir = "migrations/postgres"
	dbPort        = "5432"
	dbDriverName  = "postgres"
	dbDialect     = "postgres"
	dbUser        = "user"
	dbPassword    = "password"
	dbName        = "user_service_db"
)

func main() {
	log.Println("Starting...")
	initMigration()
	startServer()
}

func initMigration() {
	log.Println("migrations are started")
	migration := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}

	countOfMigrations, err := migrate.Exec(db, dbDialect, migration, migrate.Up)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("migrations are finished, total count: %d", countOfMigrations)
}

func startServer() {
	log.Println("server is started")
	server := grpc.NewServer()
	grpcServer := &user_service.GRPCServer{}
	api.RegisterUserServiceServer(server, grpcServer)
	listener, err := net.Listen(network, serverPort)
	if err != nil {
		log.Fatalln(err)
	}
	if err := server.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
