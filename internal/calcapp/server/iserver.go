package server

import "github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store"

type Interface interface {
	Store() store.Interface
}
