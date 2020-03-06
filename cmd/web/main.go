package main

import (
	"database/sql"
	"time"

	"github.com/allegro/bigcache"
	"github.com/brymck/helpers/cloudsqlproxy"
	"github.com/brymck/helpers/servers"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	sec "github.com/brymck/securities-service/genproto/brymck/securities/v1"
	"github.com/brymck/securities-service/pkg/models"
	"github.com/brymck/securities-service/pkg/models/mysql"
)

type application struct {
	cache  *bigcache.BigCache
	db     *sql.DB
	prices interface {
		GetMany(*time.Time, *time.Time, uint64, uint32) ([]*models.Price, error)
		Insert(*time.Time, uint64, uint32, float64) error
	}
	securities interface {
		Get(uint64) (*models.Security, error)
		Insert(string, string) (uint64, error)
	}
}

func main() {
	db, err := cloudsqlproxy.NewConnectionPool(nil)
	if err != nil {
		log.Fatal(err)
	}
	cache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}

	app := &application{
		cache:      cache,
		db:         db,
		prices:     &mysql.PriceModel{DB: db},
		securities: &mysql.SecurityModel{DB: db},
	}

	s := servers.NewGrpcServer()
	sec.RegisterSecuritiesAPIServer(s.Server, app)
	s.Serve()
}
