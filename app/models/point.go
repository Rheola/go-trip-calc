package models

import (
	"errors"
	"fmt"
)

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func (point Point) ToString() string {
	return fmt.Sprintf("(%.6f, %.6f)", point.Lat, point.Lon)
}

type RouteParams struct {
	Id       uint
	From     Point
	To       Point
	Status   uint
	Distance uint
	Duration uint
}

type CalcResult struct {
	Status   uint `json:"-"`
	Distance int  `json:"distance"`
	Duration int  `json:"duration"`
}

func (point Point) validate() error {
	if point.Lat > 90 || point.Lat < -90 {
		err := errors.New("latitude must be a number between -90 and 90")
		return err
	}

	if point.Lon > 90 || point.Lon < -90 {
		err := errors.New("longitude must be a number between -90 and 90")
		return err
	}
	return nil
}

func (params RouteParams) Validate() error {

	errFrom := params.From.validate()
	if errFrom != nil {
		err := errors.New("Wrong 'from' param: " + errFrom.Error())
		return err
	}

	errTo := params.To.validate()
	if errTo != nil {
		err := errors.New("Wrong 'to' param: " + errTo.Error())
		return err
	}

	if params.From.Lat == params.To.Lat {
		err := errors.New("From and To Lat must be difference")
		return err
	}

	if params.From.Lon == params.To.Lon {
		err := errors.New("From and To Lon must be difference")
		return err
	}
	return nil
}

// статусы проверки
const (
	StatusNone    uint = iota + 1 // 1
	StatusProcess                 // 2
	StatusOk                      // 3
	StatusFailed                  // 4
)
