package main

import (
	"encoding/json"
	"fmt"
	"go-trip-calc/app/models"
	"go-trip-calc/app/restapi"
	"io/ioutil"
	"net/http"
	"time"
)

func addRequest(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	requestModel := &models.RouteParams{}
	err = json.Unmarshal(body, requestModel)

	if err != nil {
		apiResponse := restapi.APIResponse{
			Code:    http.StatusBadRequest,
			Message: "Couldn't parse request body",
		}
		json.NewEncoder(w).Encode(apiResponse)

		return
	}

	fmt.Println(requestModel)

	// Data validation
	err = requestModel.Validate()
	if err != nil {
		fmt.Println(err)

		resp := restapi.APIResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp := restapi.APIResponse{
		Code: http.StatusCreated,
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func main() {
	//interrupt := make(chan os.Signal, 1)
	//signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	fmt.Println("Start")

	mux := http.NewServeMux()

	mux.HandleFunc("/", addRequest)

	server := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fmt.Println("starting server at :8080")

	server.ListenAndServe()
}
