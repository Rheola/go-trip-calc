package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rheola/go-trip-calc/app/models"
	"io/ioutil"
	"net/http"
)

type TravelSummary struct {
	BaseDuration int
	Duration     uint
	Length       uint
}
type section struct {
	TravelSummary TravelSummary `json:"travelSummary"`
}
type route struct {
	Sections []section `json:"sections"`
}
type answer struct {
	Routes []route `json:"routes"`
}

func CalcRoute(params models.RouteParams) (*TravelSummary, error) {

	req, err := http.NewRequest("GET", "https://router.hereapi.com/v8/routes", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("transportMode", "car")
	q.Add("origin", fmt.Sprintf("%.6f,%.6f", params.From.Lat, params.From.Lon))
	q.Add("destination", fmt.Sprintf("%.6f,%.6f", params.To.Lat, params.To.Lon))
	q.Add("return", "travelSummary")
	q.Add("apikey", "WOp3bI0eN2_gEG1ob-orRSViXwd-53mYAa_Vn8dyuMM")
	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	jsonBlob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var answer answer
	err = json.Unmarshal(jsonBlob, &answer)
	if err != nil {
		return nil, err
	}
	for _, route := range answer.Routes {
		for _, section := range route.Sections {
			return &section.TravelSummary, nil
		}
	}
	return nil, errors.New("TravelSummary not found")
}
