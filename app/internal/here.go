package internal

import (
	"encoding/json"
	"fmt"
	"github.com/rheola/go-trip-calc/app/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type travelSummary struct {
	duration int
	length   int
}

type TravelSummary struct {
	BaseDuration int
	Duration     int
	Length       int
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

func CalcRoute(params models.RouteParams) {

	req, err := http.NewRequest("GET", "https://router.hereapi.com/v8/routes", nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	q := req.URL.Query()
	q.Add("transportMode", "car")
	q.Add("origin", fmt.Sprintf("%.6f,%.6f", params.From.Lat, params.From.Lon))
	q.Add("destination", fmt.Sprintf("%.6f,%.6f", params.To.Lat, params.To.Lon))
	q.Add("return", "travelSummary")
	q.Add("apikey", "WOp3bI0eN2_gEG1ob-orRSViXwd-53mYAa_Vn8dyuMM")
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())

	// https://router.hereapi.com/v8/routes?
	//transportMode=car&
	//return=travelSummary&
	//apikey=WOp3bI0eN2_gEG1ob-orRSViXwd-53mYAa_Vn8dyuMM https://router.hereapi.com/v8/routes?transportMode=car&origin=52.5308,13.3847&destination=52.5323,13.3789&return=polyline,turnbyturnactions&apikey=WOp3bI0eN2_gEG1ob-orRSViXwd-53mYAa_Vn8dyuMM

	resp, err := http.Get(req.URL.String())
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	jsonBlob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var answer answer
	err = json.Unmarshal(jsonBlob, &answer)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(answer.Routes)
	for _, route := range answer.Routes {
		fmt.Println(route.Sections)
		for _, section := range route.Sections {
			fmt.Println(section.TravelSummary.Duration)
			fmt.Println(section.TravelSummary.BaseDuration)
			fmt.Println(section.TravelSummary.Length)
		}
	}
}
