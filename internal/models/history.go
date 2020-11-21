package models

type History struct {
	ID         int64  `json:"id" db:"id"`
	EventTime  string `json:"event_time" db:"event_time"`
	Expression string `json:"expression" db:"expression"`
	Result     string `json:"result" db:"result"`
}
