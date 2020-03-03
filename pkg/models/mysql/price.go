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

func (m *PriceModel) Insert(date *time.Time, securityId int32, priceSourceId int32, price float64) error {
	_, err := m.DB.Exec(InsertPrice, date, securityId, priceSourceId, price)
	if err != nil {
		return err
	}

	return nil
}
