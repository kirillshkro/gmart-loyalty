package types

type WithdrawRequest struct {
	Order string `json:"order"`
	Sum   int    `json:"sum"`
}
