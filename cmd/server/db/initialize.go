package db

import (
	"github.com/fev0ks/UserServiceSC/pkg/service/postgres"
	migrate "github.com/rubenv/sql-migrate"
	"log"
)

const (
	migrationsDir = "migrations/postgres"
	dbDialect     = "postgres"
)

func InitDataBase() {
	log.Println("migrations are started")
	migration := &migrate.FileMigrationSource{
		Dir: migrationsDir,
	}
	dbConnection := postgres.OpenDataBaseConnection()
	countOfMigrations, err := migrate.Exec(dbConnection, dbDialect, migration, migrate.Up)
	if err != nil {
		log.Fatalln(err)
	}
	postgres.StorageInstance = postgres.NewStorage(dbConnection)
	log.Printf("migrations are finished, total count: %d", countOfMigrations)
}
