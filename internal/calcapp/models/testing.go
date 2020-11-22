package models

import (
	"testing"
	"time"
)

func TestHistory(t *testing.T) *History {
	t.Helper()

	return &History{
		EventTime:  time.Now(),
		Expression: "(1 + 3) / 2",
		Result:     2,
	}
}
