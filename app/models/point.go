package models

import (
	"gorm.io/gorm"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type RouteParams struct {
	From Point
	To   Point
}

type routeParamsDb struct {
	gorm.Model
	From Point
	To   Point
}
