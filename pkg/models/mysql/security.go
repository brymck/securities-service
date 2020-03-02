package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/brymck/securities-service/pkg/models"
)

type SecurityModel struct {
	DB *sql.DB
}

const (
	GetSecurity  = `SELECT id, symbol, name FROM securities WHERE id = ?`
	GetBestPrice = `SELECT price FROM prices WHERE date = ? AND security_id = ? AND price_source_id = 1`
)

var (
	easternTime *time.Location
)

func init() {
	var err error
	easternTime, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(fmt.Sprintf("unable to load Eastern Time information: %v", err))
	}
}

func latestBusinessEndOfDay() time.Time {
	return time.Now().In(easternTime)
}

func (m *SecurityModel) Get(id int) (*models.Security, error) {
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

	date := "2020-02-28"
	row = m.DB.QueryRow(GetBestPrice, date, id)
	err = row.Scan(&s.Price)
	if err != nil {
		// Ignore if price is missing
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	return s, nil
}
