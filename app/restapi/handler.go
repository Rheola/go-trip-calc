package restapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rheola/go-trip-calc/app/models"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Handler struct {
	DB *models.TripDb
	Ch chan models.RouteParams
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {

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

	err = h.DB.Create(requestModel)

	if err != nil {
		ResponseInternalError(w, err)
		return
	}
	h.Ch <- *requestModel
	resp := APIResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("%d", requestModel.Id),
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		ResponseInternalError(w, errors.New("Id not set"))
		return
	}

	calcResult, err := h.DB.Get(id)
	if err != nil {
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
