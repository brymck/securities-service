package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/brymck/securities-service/pkg/dates"
	"github.com/brymck/securities-service/pkg/models"
)

type SecurityModel struct {
	DB *sql.DB
}

const (
	GetSecurity    = `SELECT id, symbol, name FROM securities WHERE id = ?`
	GetBestPrice   = `SELECT price FROM prices WHERE date = ? AND security_id = ? AND price_source_id = 1`
	InsertSecurity = `INSERT INTO securities (symbol, name) VALUES (?, ?)`
)

var (
	easternTime *time.Location
	now         = time.Now
)

func init() {
	var err error
	easternTime, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Sprintf("unable to load Eastern Time information: %v", err))
	}
}

func (m *SecurityModel) Get(id int32) (*models.Security, error) {
	s := &models.Security{}

	row := m.DB.QueryRow(GetSecurity, id)
	err := row.Scan(&s.ID, &s.Symbol, &s.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	date := dates.LatestBusinessEndOfDay(now().In(easternTime))
	row = m.DB.QueryRow(GetBestPrice, dates.IsoDate(date), id)
	err = row.Scan(&s.Price)
	if err != nil {
		// Ignore if price is missing
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return s, nil
}

func (m *SecurityModel) Insert(symbol string, name string) (int32, error) {
	r, err := m.DB.Exec(InsertSecurity, symbol, name)
	if err != nil {
		return 0, err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}
