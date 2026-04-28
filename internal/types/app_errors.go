package types

import "fmt"

type ErrDuplicateUser struct {
	UserName string
}

func (e ErrDuplicateUser) Error() string {
	return fmt.Sprintf("user with name %s already exists", e.UserName)
}
