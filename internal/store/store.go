package store

type Store interface {
	History() HistoryRepository
}
