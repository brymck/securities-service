package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/bmizerany/pat"
	"github.com/brymck/helpers/webapp"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"

	"github.com/brymck/securities-service/pkg/models"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
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
		price, err := getPrice(s.Symbol)
		if err != nil {
			webapp.ServerError(w, err)
		}
		s.Price = price
	}

	err = json.NewEncoder(w).Encode(s)
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

	ts, err := getHistoricalPrices(s.Symbol)
	if err != nil {
		webapp.ServerError(w, err)
	}

	for _, item := range ts {
		err = app.prices.Insert(item.Date, s.ID, 1, item.Price)
		if err != nil {
			webapp.ServerError(w, err)
		}
	}
}

func (app *application) Routes() http.Handler {
	standardMiddleware := alice.New(app.logRequest)

	mux := pat.New()
	mux.Get("/security/:id", standardMiddleware.ThenFunc(app.getSecurity))
	mux.Get("/security/:id/update-prices", standardMiddleware.ThenFunc(app.updatePrices))

	return standardMiddleware.Then(mux)
}
