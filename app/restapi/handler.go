package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rheola/go-trip-calc/app/internal"
	"github.com/rheola/go-trip-calc/app/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Handler struct {
	DB     *models.TripDb
	Ch     chan models.RouteParams
	Wg     sync.WaitGroup
	Closed chan struct{}
}

func (handler *Handler) Add(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Wg.Add Add")
	handler.Wg.Add(1)
	defer handler.Wg.Done()
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	requestModel := &models.RouteParams{}
	err = json.Unmarshal(body, requestModel)

	if err != nil {
		errBody := errors.New("Couldn't parse request body")
		ResponseBadRequest(w, errBody)
		return
	}

	// Data validation
	err = requestModel.Validate()
	if err != nil {
		ResponseBadRequest(w, err)
		return
	}

	err = handler.DB.Create(requestModel)

	if err != nil {
		ResponseInternalError(w, err)
		return
	}
	handler.Ch <- *requestModel
	resp := APIResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("%d", requestModel.Id),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (handler *Handler) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Wg Add Get")

	handler.Wg.Add(1)
	fmt.Println("Get start")
	defer handler.Wg.Done()

	time.Sleep(5 * time.Second)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		ResponseInternalError(w, errors.New("Id not set"))
		return
	}

	calcResult, err := handler.DB.Get(id)
	if err != nil {
		fmt.Println(err)
		if err.Error() == "sql: no rows in result set" {
			ResponseNotFoundError(w, err)
			return
		}
		ResponseInternalError(w, err)
		return
	}
	resp := APIResponse{
		Code: http.StatusOK,
	}
	switch calcResult.Status {
	case models.StatusNone:
		fallthrough
	case models.StatusProcess:
		resp.Message = "Waiting. Route not calking yet"
	case models.StatusFailed:
		resp.Message = "Calc error"
	case models.StatusOk:
		resp.Message = *calcResult
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (handler *Handler) Stop() {
	close(handler.Closed)
	fmt.Println("Stop")
	handler.Wg.Wait()
	fmt.Println("after Stop ")
}

func (handler *Handler) Run() {
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

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen:%+s\n", err)
		} else {
			fmt.Println("listen")
		}
	}()

}

func (handler *Handler) Worker(val models.RouteParams) {
	defer handler.Wg.Done()

	fmt.Println("Worker start")

	time.Sleep(10 * time.Second)

	val.Status = models.StatusProcess
	err := handler.DB.Update(&val)

	if err != nil {
		fmt.Println("Db update error")
		fmt.Println(err)
		return
	}

	res, err := internal.CalcRoute(val)

	if err != nil {
		val.Status = models.StatusFailed
		fmt.Println("CalcRoute error ", err)

		err = handler.DB.Update(&val)
		if err != nil {
			fmt.Println("Db update error", err)
		}
		return
	}

	val.Status = models.StatusOk
	val.Distance = res.Length
	val.Duration = res.Duration
	err = handler.DB.Update(&val)
	if err != nil {
		fmt.Println("Db update error", err)
	}
	fmt.Println("Worker end")
}
