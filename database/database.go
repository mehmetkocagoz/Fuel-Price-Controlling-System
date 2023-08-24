package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbname   = "brentoil"
	user     = "postgres"
	password = "123456"
	host     = "localhost"
	port     = "5432"
	sslmode  = "disable"
)

func Connect() *sql.DB {
	psglconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
	db, err := sql.Open("postgres", psglconn)
	if err != nil {
		fmt.Println("Database connection has failed: ", err)
	}
	fmt.Println("Database is connected.")
	return db
}
