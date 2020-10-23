package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rheola/go-trip-calc/app/internal"
	"github.com/rheola/go-trip-calc/app/models"
	"github.com/rheola/go-trip-calc/app/resource"
	"github.com/rheola/go-trip-calc/app/restapi"
	"net/http"
	"time"
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
		DB: &tripDb,
		Ch: hereChannel,
	}

	go func(in chan models.RouteParams) {
		for {
			select {
			case val := <-in:
				val.Status = models.StatusProcess
				err := handler.DB.Update(&val)

				if err != nil {
					fmt.Println(" Update err")
					fmt.Println(err)

				} else {

					res, err := internal.CalcRoute(val)

					if err != nil {
						val.Status = models.StatusFailed
						fmt.Println(err)
						_ = handler.DB.Update(&val)
					} else {

						val.Status = models.StatusOk
						val.Distance = res.Length
						val.Duration = res.Duration
						err = handler.DB.Update(&val)
						if err != nil {

						}
					}
				}

			}
		}
	}(hereChannel)

	//mux := http.NewServeMux()
	mux := mux.NewRouter()
	mux.HandleFunc("/routes/{id:[0-9]+}", handler.Get)
	mux.HandleFunc("/routes", handler.Add)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("starting server at :8080")

	server.ListenAndServe()

}
