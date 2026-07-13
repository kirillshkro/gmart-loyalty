package types

import "time"

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}

type WithdrawResponse struct {
	Order       string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
