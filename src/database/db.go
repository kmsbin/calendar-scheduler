package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "1234"
	dbname   = "skill_share"
)

func OpenConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = conn.Ping()
	if err != nil {
		panic(err)
	}
	return conn, err
}

func CloseConnection(conn *sql.DB) {
	if err := conn.Close(); err != nil {
		log.Fatalf("Error opening db connection %v", err)
	}
}
