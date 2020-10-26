package main

import (
	"database/sql"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rheola/go-trip-calc/app/models"
	"github.com/rheola/go-trip-calc/app/resource"
	"github.com/rheola/go-trip-calc/app/restapi"
	"os"
	"os/signal"
)

func main() {

	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	conf := resource.Config{}
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	dbConn, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		panic(err)
	}
	err = dbConn.Ping() // вот туc будет первое подключение к базе
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	tripDb := models.TripDb{*dbConn}

	hereChannel := make(chan models.RouteParams)
	handler := &restapi.Handler{
		DB:     &tripDb,
		Ch:     hereChannel,
		Closed: make(chan struct{}),
	}

	closeChannel := make(chan os.Signal)
	signal.Notify(closeChannel, os.Interrupt)

	go func(in chan models.RouteParams) {
		for {
			select {
			case val := <-in:
				handler.Wg.Add(1)
				handler.Worker(val)
			}
		}
	}(hereChannel)

	go func() {
		handler.Run()
	}()

	select {
	case sig := <-closeChannel:
		fmt.Printf("okunewa Got %s signal. Aborting...\n", sig)
		handler.Stop()
	}
}
