package models

import (
	"errors"
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
