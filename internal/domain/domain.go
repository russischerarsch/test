package domain

import "time"

type Quote struct {
	ID        int
	Pair      string
	Price     *float64
	UpdatedAt *time.Time
	Status    string
}
