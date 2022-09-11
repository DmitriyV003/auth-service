package main

import (
	"auth/data"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"log"
	"net/http"
	"os"
	"time"
)

const webPOrt = "80"

var counts int64

type Config struct {
	Db     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting auth service")

	conn := connectToDb()
	if conn == nil {
		log.Panic("Can't connect to postgres")
	}

	app := Config{
		Db:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPOrt),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDb(dsn string) (*sql.DB, error) {
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

func connectToDb() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		conn, err := openDb(dsn)
		if err != nil {
			log.Printf("Unable to connect to DB")
			counts++
		} else {
			log.Printf("Connected to Postgres")
			return conn
		}

		if counts > 10 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for 2 sec")
		time.Sleep(2 * time.Second)
		continue
	}
}
