package mysql

import (
	"database/sql"
	"time"
)

const (
	InsertPrice = `INSERT INTO prices (date, security_id, price_source_id, price)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE price = VALUES(price)`
)

type PriceModel struct {
	DB *sql.DB
}

func (m *PriceModel) Insert(date *time.Time, securityId uint64, priceSourceId uint32, price float64) error {
	_, err := m.DB.Exec(InsertPrice, date, securityId, priceSourceId, price)
	if err != nil {
		return err
	}

	return nil
}
