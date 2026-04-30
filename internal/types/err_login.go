package types

import "fmt"

type ErrInvalidLogin struct {
	UserName string
}

func (e *ErrInvalidLogin) Error() string {
	return fmt.Sprintf("invalid login for user %s", e.UserName)
}
