package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/models"
	"time"
)

type HistoryRepository struct {
	store *Store
}

func (r *HistoryRepository) Create(h *models.History) error {
	_, err := r.store.db.Exec(
		"INSERT INTO history (event_time, expression, result) VALUES ($1, $2, $3)",
		h.EventTime.Format(time.RFC3339),
		h.Expression,
		h.Result,
	)
	return err
}

func (r *HistoryRepository) GetHistoryByTimeRange(startTime string, endTime string) (history []*models.History, err error) {
	var rows *sql.Rows

	rows, err = r.store.db.Query("SELECT id, event_time, expression, result FROM history WHERE event_time BETWEEN $1 and $2", startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("sql query error: %w", err)
	}

	var eventTime string
	for rows.Next() {
		h := new(models.History)
		err = rows.Scan(&h.ID, &eventTime, &h.Expression, &h.Result)
		if err != nil {
			return nil, fmt.Errorf("rows scanning error: %w", err)
		}

		tim, err := time.Parse(time.RFC3339, eventTime)
		if err != nil {
			return nil, fmt.Errorf("event time parsing error: %w", err)
		}
		h.EventTime = tim
		history = append(history, h)
	}
	return history, nil
}

func (r *HistoryRepository) GetAllHistory() (history []*models.History, err error) {
	var rows *sql.Rows

	rows, err = r.store.db.Query("SELECT id, event_time, expression, result FROM history")
	if err != nil {
		return nil, fmt.Errorf("sql query error: %w", err)
	}

	var eventTime string
	for rows.Next() {
		h := new(models.History)
		err = rows.Scan(&h.ID, &eventTime, &h.Expression, &h.Result)
		if err != nil {
			return nil, fmt.Errorf("rows scanning error: %w", err)
		}

		tim, err := time.Parse(time.RFC3339, eventTime)
		if err != nil {
			return nil, fmt.Errorf("event time parsing error: %w", err)
		}
		h.EventTime = tim
		history = append(history, h)
	}
	return history, nil
}
