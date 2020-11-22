package server

import (
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store"
	"github.com/sirupsen/logrus"
)

type Interface interface {
	Store() store.Interface
	Logger() *logrus.Logger
}
