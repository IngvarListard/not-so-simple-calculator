package models

import "time"

type History struct {
	ID         int64       `json:"id" db:"id"`
	EventTime  time.Time   `json:"event_time" db:"event_time"`
	Expression string      `json:"expression" db:"expression"`
	Result     interface{} `json:"result" db:"result"`
}
