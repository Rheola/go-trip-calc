package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rheola/go-trip-calc/app/models"
	"github.com/rheola/go-trip-calc/app/restapi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func handleRequests() {
	router := mux.NewRouter()

	router.HandleFunc("/calc", addRequest()).Methods(http.MethodPost)
	//router.HandleFunc("/status", checkStatus()).Methods(http.MethodGet)

	fmt.Println("PORT:", os.Getenv("PORT"))
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, router))
}

func addRequest() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		routeParams := new(models.RouteParams)
		err := json.NewDecoder(r.Body).Decode(routeParams)
		if err != nil {
			apiResponse := restapi.APIResponse{
				Code:    http.StatusBadRequest,
				Message: "Couldn't parse request body",
			}
			json.NewEncoder(w).Encode(apiResponse)

			return
		}

		// Data validation

		// Pet creation
		/*h.petController.AddPet(petItem)

		resp := APIResponse{
			Code: http.StatusCreated,
		}*/
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}
}

func checkStatus(w http.ResponseWriter, r *http.Request) {

}

func createNewArticle(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	//db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	fmt.Println("Start")

	handleRequests()
}
