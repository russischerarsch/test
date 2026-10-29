package domain

import "time"

type Quote struct {
	ID        int
	Pair      string `json:"pair"`
	Price     float64
	UpdatedAt time.Time
	Status    string
}
