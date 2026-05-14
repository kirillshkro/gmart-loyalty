package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// genNumberOrder генерирует случайную строку длины length, представляющую валидное по алгоритму Луна число.
func GenNumberOrder(length int) string {
	if length < 2 {
		return ""
	}

	// Генерируем все цифры, кроме последней
	digits := make([]string, length-1)
	for i := range digits {
		digits[i] = fmt.Sprintf("%d", rand.Intn(10))
	}

	// Вычисляем контрольную цифру по алгоритму Луна
	sum := 0
	multiplier := 2
	for i := len(digits) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(digits[i])
		if multiplier == 2 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		multiplier = 3 - multiplier
	}

	// Вычисляем контрольную цифру
	checksum := sum % 10
	checkdigit := (10 - checksum) % 10

	// Добавляем контрольную цифру в конец строки
	digits = append(digits, fmt.Sprintf("%d", checkdigit))

	return strings.Join(digits, "")
}
