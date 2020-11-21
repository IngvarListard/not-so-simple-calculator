package store

type Interface interface {
	History() HistoryRepository
}
