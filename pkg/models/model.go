package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")
)

type Security struct {
	ID     uint64
	Symbol string
	Name   string
	Price  float64
}

type Price struct {
	Date          *time.Time
	SecurityId    uint64
	PriceSourceId uint32
	Price         float64
}
