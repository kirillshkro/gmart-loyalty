package utils

import (
	"strconv"
	"strings"
)

// Проверка ID на валидность, по алгоритму Луна.
func Valid(orderNum string) bool {
	orderNum = strings.ReplaceAll(orderNum, " ", "")
	orderNum = strings.ReplaceAll(orderNum, "-", "")

	if orderNum == "" {
		return false
	}

	sum := 0
	isSecond := false

	for i := len(orderNum) - 1; i >= 0; i-- {
		dig, err := strconv.Atoi(string(orderNum[i]))
		if err != nil {
			return false
		}

		if isSecond {
			dig *= 2
			if dig > 9 {
				dig -= 9
			}
		}
		sum += dig
		isSecond = !isSecond
	}

	return sum%10 == 0
}
