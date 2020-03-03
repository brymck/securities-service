package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)


func (app *application) Routes() http.Handler {
	standardMiddleware := alice.New(app.logRequest)

	mux := pat.New()
	mux.Get("/security/:id", standardMiddleware.ThenFunc(app.getSecurity))
	mux.Get("/security/:id/update-prices", standardMiddleware.ThenFunc(app.updatePrices))
	mux.Post("/security/", standardMiddleware.ThenFunc(app.insertSecurity))

	return standardMiddleware.Then(mux)
}
