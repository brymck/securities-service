package main

import (
	"database/sql"
	"time"

	"github.com/brymck/helpers/cloudsqlproxy"
	"github.com/brymck/helpers/servers"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"

	pb "github.com/brymck/securities-service/genproto"
	"github.com/brymck/securities-service/pkg/models"
	"github.com/brymck/securities-service/pkg/models/mysql"
)

type application struct {
	db     *sql.DB
	prices interface {
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

	app := &application{
		db:         db,
		prices:     &mysql.PriceModel{DB: db},
		securities: &mysql.SecurityModel{DB: db},
	}

	s := servers.NewGrpcServer()
	pb.RegisterSecuritiesAPIServer(s.Server, app)
	s.Serve()
}
