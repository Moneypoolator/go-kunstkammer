package utils

import (
	"fmt"
	"regexp"
	"strings"
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

// ExtractWorkCode extracts the product code (e.g., "CAD", "MGM") and the work code from a job title string.
// The work code returned is the part after the first dot (e.g., for 'US.19.03', returns '19.03').
// The first return value is the product code.
// The second return value is the work code (after the first dot).
// The third return value is the error.
//
// The expected format is: [PRODUCT]:work.code.part
//
// Example usage:
//
//	product, workCode, err := ExtractWorkCode("[CAD]:TS.FEATURE.123. Task Title")
//	// product == "CAD", workCode == "FEATURE.123", err == nil
//
//	product, workCode, err := ExtractWorkCode("[MGM]:US.19.03. Some Task")
//	// product == "MGM", workCode == "19.03", err == nil
//
//	product, workCode, err := ExtractWorkCode("NoCodeHere")
//	// product == "", workCode == "", err != nil
func ExtractWorkCode(parentTitle string) (string, string, error) {
	re := regexp.MustCompile(`\[([A-Za-z]+)\]:([A-Za-z]+\.[^\.\s]+\.[^\.\s]+)`) // Matches [PRODUCT]:work.code.part
	match := re.FindStringSubmatch(parentTitle)

	if len(match) < 3 {
		return "", "", fmt.Errorf("work code not found in title: %s", parentTitle)
	}

	// match[2] is like 'US.19.03' or 'TS.FEATURE.123'
	workCodeFull := match[2]
	dotIdx := strings.Index(workCodeFull, ".")
	if dotIdx == -1 || dotIdx+1 >= len(workCodeFull) {
		return match[1], "", fmt.Errorf("work code format invalid in title: %s", parentTitle)
	}
	workCode := workCodeFull[dotIdx+1:]

	return match[1], workCode, nil
}
