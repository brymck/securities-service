package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/brymck/helpers/webapp"
	log "github.com/sirupsen/logrus"

	pb "github.com/brymck/securities-service/genproto"
	"github.com/brymck/securities-service/pkg/models"
)

type InsertSecurityRequest struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func (app *application) getSecurity(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		webapp.NotFound(w)
		return
	}

	s, err := app.securities.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			webapp.NotFound(w)
		} else {
			webapp.ServerError(w, err)
		}
		return
	}

	if s.Price == 0.0 {
		log.Infof("retrieving missing price for %s", s.Symbol)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := alphaVantageApi.GetQuote(ctx, &pb.GetQuoteRequest{Symbol: s.Symbol})
		if err != nil {
			webapp.ServerError(w, err)
			return
		}
		s.Price = resp.Price
	}

	err = json.NewEncoder(w).Encode(s)
	if err != nil {
		webapp.ServerError(w, err)
	}
}

func (app *application) insertSecurity(w http.ResponseWriter, r *http.Request) {
	v := &InsertSecurityRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		webapp.ClientError(w, http.StatusBadRequest)
		return
	}

	if v.Symbol == "" {
		log.Error("symbol cannot be blank")
		webapp.ClientError(w, http.StatusBadRequest)
		return
	}

	if v.Name == "" {
		log.Error("name cannot be blank")
		webapp.ClientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.securities.Insert(v.Symbol, v.Name)
	if err != nil {
		webapp.ServerError(w, err)
		return
	}

	err = json.NewEncoder(w).Encode(&models.Security{ID: id})
	if err != nil {
		webapp.ServerError(w, err)
	}
}

func (app *application) updatePrices(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		webapp.NotFound(w)
		return
	}

	s, err := app.securities.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			webapp.NotFound(w)
		} else {
			webapp.ServerError(w, err)
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	resp, err := alphaVantageApi.GetTimeSeries(ctx, &pb.GetTimeSeriesRequest{Symbol: s.Symbol, Full: true})
	if err != nil {
		webapp.ServerError(w, err)
	}

	ts := make([]*DatePrice, len(resp.TimeSeries))
	for i, item := range resp.TimeSeries {
		date := time.Date(int(item.Date.Year), time.Month(item.Date.Month), int(item.Date.Day), 0, 0, 0, 0, time.UTC)
		ts[i] = &DatePrice{Date: &date, Price: item.Close}
	}

	start := time.Now()
	log.Infof("inserting %d historical prices", len(ts))
	for _, item := range ts {
		err = app.prices.Insert(item.Date, s.ID, 1, item.Price)
		if err != nil {
			webapp.ServerError(w, err)
		}
	}
	end := time.Now()
	log.Infof("inserted %d historical prices in %d ms", len(ts), end.Sub(start).Milliseconds())
}
