package utils

import (
	"fmt"
	"regexp"
)

// Вспомогательные функции для создания указателей
func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func Ptr[T any](value T) *T {
	return &value
}

func ExtractWorkCode(parentTitle string) (string, error) {
	re := regexp.MustCompile(`\[CAD\]:[A-Za-z]+\.([^.\s]+\.[^.\s]+)`)
	match := re.FindStringSubmatch(parentTitle)

	if len(match) < 2 {
		return "", fmt.Errorf("work code not found in title: %s", parentTitle)
	}

	return match[1], nil
}
