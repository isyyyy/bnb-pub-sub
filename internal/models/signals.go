package models

import "time"

type Signal struct {
	Symbol    string    `json:"symbol"`
	Entry     string    `json:"entry"`
	Target    string    `json:"target"`
	StartTime time.Time `json:"startTime"`
	Duration  string    `json:"duration"`
}
