package types

import "fmt"

type ErrInvalidFormatOrder struct {
	Number string
}

func (e *ErrInvalidFormatOrder) Error() string {
	return fmt.Sprintf("invalid format order number: %s", e.Number)
}
