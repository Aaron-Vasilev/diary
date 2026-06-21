package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

var Query *Queries

func ConnectDB() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Error connecting to db: ", err)
	}
	Query = New(conn)
	return conn
}
