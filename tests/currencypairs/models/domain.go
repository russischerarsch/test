package models

import "time"

type CreatePairReq struct {
	Pair string `json:"Pair"`
}

type CreatePairResp struct {
	ID string `json:"id"`
}

type Quote struct {
	ID        string `json:"id"`
	Pair      string `json:"pair"`
	Price     *float64
	UpdatedAt *time.Time
	CreatedAt time.Time
	Status    string
}

type ErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
