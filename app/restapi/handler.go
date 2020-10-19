package restapi

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rheola/go-trip-calc/app/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	DB *sql.DB
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

	var id int
	now := time.Now()
	err = h.DB.QueryRow(
		"INSERT INTO rates (from_point, to_point, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $4) RETURNING id",
		requestModel.From.ToString(),
		requestModel.To.ToString(),
		models.StatusNone,
		now.Format("2006-01-02 15:04:05"),
	).Scan(&id)

	if err != nil {
		ResponseInternalError(w, err)
		return
	}
	h.Ch <- *requestModel
	resp := APIResponse{
		Code:    http.StatusCreated,
		Message: fmt.Sprintf("%d", id),
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

	calcResult := &models.CalcResult{}
	// QueryRow сам закрывает коннект
	row := h.DB.QueryRow("SELECT status, distance, duration FROM rates  WHERE id = $1", id)
	err = row.Scan(&calcResult.Status, &calcResult.Distance, &calcResult.Duration)
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
