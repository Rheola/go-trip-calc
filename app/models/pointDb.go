package models

import (
	"database/sql"
	"time"
)

type TripDb struct {
	sql.DB
}

func (db *TripDb) Update(requestModel *RouteParams) error {
	now := time.Now()
	_, err := db.Exec(
		"Update rates set status = $2, distance = $3, duration=$4,  updated_at= $5  where id =  $1 ",
		requestModel.Id,
		requestModel.Status,
		requestModel.Distance,
		requestModel.Duration,
		now.Format("2006-01-02 15:04:05"),
	)
	return err
}

func (db *TripDb) Create(requestModel *RouteParams) error {
	var id uint
	now := time.Now()
	err := db.QueryRow(
		"INSERT INTO rates (from_point, to_point, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $4) RETURNING id",
		requestModel.From.ToString(),
		requestModel.To.ToString(),
		StatusNone,
		now.Format("2006-01-02 15:04:05"),
	).Scan(&id)
	requestModel.Id = id
	return err
}

func (db *TripDb) Get(id int) (*CalcResult, error) {
	calcResult := &CalcResult{}
	row := db.QueryRow("SELECT status, distance, duration FROM rates  WHERE id = $1", id)
	err := row.Scan(&calcResult.Status, &calcResult.Distance, &calcResult.Duration)
	if err != nil {
		return nil, err
	}
	return calcResult, nil
}
