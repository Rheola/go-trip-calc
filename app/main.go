package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rheola/go-trip-calc/app/config"
	"github.com/rheola/go-trip-calc/app/models"
	"github.com/rheola/go-trip-calc/app/restapi"
	"log"
	"os"
	"os/signal"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()
	fmt.Println(conf.Port)

	fmt.Println(conf.DbUrl)
	dbConn, err := sql.Open("postgres", conf.DbUrl)

	if err != nil {
		panic(err)
	}

	err = dbConn.Ping()
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	tripDb := models.TripDb{*dbConn}

	hereChannel := make(chan models.RouteParams, 1)

	handler := &restapi.Handler{
		DB:     &tripDb,
		Ch:     hereChannel,
		Closed: make(chan struct{}),
	}

	closeChannel := make(chan os.Signal)
	signal.Notify(closeChannel, os.Interrupt)

	go func() {
		handler.Run()
	}()

	select {
	case _ = <-closeChannel:
		handler.Stop()
	}
}
