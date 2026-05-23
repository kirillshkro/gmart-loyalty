package types

import "fmt"

type ErrDuplicateUser struct {
	UserName string
}

func (e ErrDuplicateUser) Error() string {
	return fmt.Sprintf("user with name %s already exists", e.UserName)
}

type ErrInsufficientBalance struct {
	UserID  int
	Balance float64
}

func (e ErrInsufficientBalance) Error() string {
	return fmt.Sprintf("user %d has insufficient balance for requested operation. Balance: %.2f", e.UserID, e.Balance)
}
