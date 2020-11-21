package sqlstore

import (
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/internal/models"
	"time"
)

type HistoryRepository struct {
	store *Store
}

func (r *HistoryRepository) Create(h *models.History) error {
	return r.store.db.QueryRow(
		"INSERT INTO history (event_time, expresstion, result) VALUES ($1, $2, $3) RETURNING id",
		h.EventTime,
		h.Expression,
		h.Result,
	).Scan(&h.ID)
}

func (r *HistoryRepository) GetHistoryFromTime(startTime *time.Time) ([]*models.History, error) {
	var history []*models.History
	rows, err := r.store.db.Query("SELECT id, event_time, expression, result FROM history WHERE event_time BETWEEN $1 and now()", startTime)
	if err != nil {
		return nil, fmt.Errorf("sql query error: %w", err)
	}

	for rows.Next() {
		h := new(models.History)
		err = rows.Scan(&h.ID, &h.EventTime, &h.Expression, &h.Result)
		if err != nil {
			return nil, fmt.Errorf("rows scanning error: %w", err)
		}
		history = append(history, h)
	}
	return history, nil
}
