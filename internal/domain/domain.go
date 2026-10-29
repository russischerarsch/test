package domain

import "time"

type Quote struct {
	ID        int
	Pair      string `json:"pair"`
	Price     *float64
	UpdatedAt *time.Time
	CreatedAt time.Time
	Status    string
}

type Info struct {
	Success bool               `json:"success"`
	Rates   map[string]float64 `json:"rates"`
	Error   *struct {
		Code int    `json:"code"`
		Info string `json:"info"`
	} `json:"error"`
}
