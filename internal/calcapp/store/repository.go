package store

import (
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/models"
)

type HistoryRepository interface {
	Create(history *models.History) error
	GetHistoryByTimeRange(startTime string, endTime string) ([]*models.History, error)
	GetAllHistory() ([]*models.History, error)
}
