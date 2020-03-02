package main

import (
	"database/sql"
	"time"

	"github.com/brymck/helpers/cloudsqlproxy"
	"github.com/brymck/helpers/webapp"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	"github.com/brymck/securities-service/pkg/models"
	"github.com/brymck/securities-service/pkg/models/mysql"
)

type application struct {
	db     *sql.DB
	prices interface {
		Insert(*time.Time, int, int, float64) error
	}
	securities interface {
		Get(int) (*models.Security, error)
	}
}

func main() {
	db, err := cloudsqlproxy.NewConnectionPool(nil)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		db:         db,
		prices:     &mysql.PriceModel{DB: db},
		securities: &mysql.SecurityModel{DB: db},
	}

	webapp.Serve(app)
}
