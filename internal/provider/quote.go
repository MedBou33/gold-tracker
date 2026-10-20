package provider

import "time"

type Quote struct {
	Currency  string    `json:"currency"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	Symbol    string    `json:"symbol"`
	UpdatedAt time.Time `json:"updatedAt"`
}
