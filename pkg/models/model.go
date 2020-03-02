package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Security struct {
	ID     int
	Symbol string
	Name   string
	Price  float64
}

type Price struct {
	Date          *time.Time
	SecurityId    int
	PriceSourceId int
	Price         float64
}
