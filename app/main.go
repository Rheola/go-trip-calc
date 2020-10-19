package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"github.com/rheola/go-trip-calc/app/models"
	"github.com/rheola/go-trip-calc/app/restapi"
	"io/ioutil"
	"net/http"
	"time"
)

type Handler struct {
	DB *sql.DB
}

func (h *Handler) add(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	requestModel := &models.RouteParams{}
	err = json.Unmarshal(body, requestModel)

	if err != nil {
		err1 := errors.New("Couldn't parse request body")
		ResponseBadRequest(w, err1)
		return
	}

	// Data validation
	err = requestModel.Validate()
	if err != nil {
		ResponseBadRequest(w, err)
		return
	}

	var id int
	err = h.DB.QueryRow(
		"INSERT INTO rates (from_point, to_point) VALUES ($1, $2) RETURNING id",
		requestModel.From.ToString(),
		requestModel.To.ToString(),
	).Scan(&id)

	if err != nil {
		ResponseInternalError(w, err)
		return
	}

	resp := restapi.APIResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("%d", id),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

func ResponseBadRequest(w http.ResponseWriter, err error) {
	fmt.Println(err)

	resp := restapi.APIResponse{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func ResponseInternalError(w http.ResponseWriter, err error) {
	fmt.Println(err)

	resp := restapi.APIResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func main() {
	conf := Config{}
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	fmt.Println(conf.DBURL)

	dbConn, err := sql.Open("postgres", conf.DBURL)
	if err != nil {
		panic(err)
	}
	//db.SetMaxOpenConns(10)
	err = dbConn.Ping() // вот тут будет первое подключение к базе
	if err != nil {
		panic(err)
	}
	defer dbConn.Close()

	handler := &Handler{
		DB: dbConn,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.add)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("starting server at :8080")

	server.ListenAndServe()
}

type Config struct {
	//dsn := "root@tcp(localhost:3306)/coursera?"
	// указываем кодировку
	//dsn += "&charset=utf8"
	// отказываемся от prapared statements
	// параметры подставляются сразу
	//dsn += "&interpolateParams=true"
	//RESTAPIPort int    `envconfig:"PORT" default:"8080" required:"true"`
	DBURL string `envconfig:"DB_URL" default:"postgres://db-user:db-password@localhost:5429/tripdb?sslmode=disable" required:"true"`
}
