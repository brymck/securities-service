package mysql

import (
	"database/sql"
	"time"

	"github.com/brymck/securities-service/pkg/models"
)

const (
	GetManyPrices = `SELECT date, price
		FROM prices
		WHERE date BETWEEN ? AND ?
			AND security_id = ?
			AND price_source_id = ?`
	InsertPrice = `INSERT INTO prices (date, security_id, price_source_id, price)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE price = VALUES(price)`
)

type PriceModel struct {
	DB *sql.DB
}

func (m *PriceModel) GetMany(startDate *time.Time, endDate *time.Time, securityId uint64, priceSourceId uint32) ([]*models.Price, error) {
	rows, err := m.DB.Query(GetManyPrices, startDate, endDate, securityId, priceSourceId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var prices []*models.Price
	for rows.Next() {
		price := &models.Price{SecurityId: securityId, PriceSourceId: priceSourceId}
		err = rows.Scan(&price.Date, &price.Price)
		if err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return prices, nil
}

func (m *PriceModel) Insert(date *time.Time, securityId uint64, priceSourceId uint32, price float64) error {
	_, err := m.DB.Exec(InsertPrice, date, securityId, priceSourceId, price)
	if err != nil {
		return err
	}

	return nil
}
