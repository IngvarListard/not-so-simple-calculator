package store

import (
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/models"
	"time"
)

type HistoryRepository interface {
	Create(history *models.History) error
	GetHistoryFromTime(startTime *time.Time) ([]*models.History, error)
}
