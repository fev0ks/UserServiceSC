package postgres

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	host         = "localhost"
	dbPort       = "5432"
	dbDriverName = "postgres"
	dbUser       = "user"
	dbPassword   = "password"
	dbName       = "user_service_db"
)

func OpenDataBaseConnection() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open(dbDriverName, psqlInfo)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}

func CloseDataBaseConnection(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Fatalln(err)
	}
}
