package types

import "fmt"

type ErrOwnAnotherUser struct {
	UserID   int
	OrderNum string
}

func (e ErrOwnAnotherUser) Error() string {
	return fmt.Sprintf("user %d can't own order %s", e.UserID, e.OrderNum)
}
