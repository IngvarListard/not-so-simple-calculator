package sqlstore

import (
	"database/sql"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store"
)

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

type Store struct {
	db                *sql.DB
	historyRepository *HistoryRepository
}

func (s *Store) History() store.HistoryRepository {
	if s.historyRepository == nil {
		s.historyRepository = &HistoryRepository{s}
	}
	return s.historyRepository
}
