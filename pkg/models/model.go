package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Security struct {
	ID     int32
	Symbol string
	Name   string
	Price  float64
}

type Price struct {
	Date          *time.Time
	SecurityId    int32
	PriceSourceId int32
	Price         float64
}
