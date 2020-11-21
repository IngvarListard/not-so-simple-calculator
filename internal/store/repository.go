package store

import (
	"github.com/IngvarListard/not-so-simple-calculator/internal/models"
	"time"
)

type HistoryRepository interface {
	Create(history *models.History) error
	GetHistoryFromTime(startTime *time.Time) ([]*models.History, error)
}
