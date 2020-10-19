package restapi

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}

func ResponseBadRequest(w http.ResponseWriter, err error) {

	resp := APIResponse{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}

func ResponseInternalError(w http.ResponseWriter, err error) {

	resp := APIResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(resp)
}
