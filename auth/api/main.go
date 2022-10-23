package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/radugaf/RelantinV3/auth/models"
)

const port = "80"

var counts int64

type Config struct {
	DB     *sql.DB
	Models models.Models
}

func main() {
	log.Println("Starting authentication service")

	// Connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to database.")
	}

	// Set up config
	app := Config{
		DB:     conn,
		Models: models.New(conn),
	}

	service := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := service.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}

}

// 'dsn' stands for data source name
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to database")
			return conn
		}

		if counts > 10 {
			log.Println("Could not connect to database. Tried 10 times. Exiting...")
			return nil
		}

		log.Println("Waiting 5 seconds before trying again")
		time.Sleep(5 * time.Second)

		continue
	}
}
