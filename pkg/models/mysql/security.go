package mysql

import (
	"database/sql"
	"errors"

	"github.com/brymck/securities-service/pkg/models"
)

type SecurityModel struct {
	DB *sql.DB
}

func (m *SecurityModel) Get(id int) (*models.Security, error) {
	stmt := `SELECT id, symbol, name FROM securities WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Security{}

	err := row.Scan(&s.ID, &s.Symbol, &s.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}
