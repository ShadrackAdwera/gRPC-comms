package main

import (
	"categories/repo"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	webPort = "5000"
)

var counts int

type Config struct {
	DB     *sql.DB
	Models repo.Models
}

func main() {

	conn := connectToDB()

	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	app := Config{
		DB:     conn,
		Models: repo.New(conn),
	}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.Router(),
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}
