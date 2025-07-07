package snake

import (
	"strings"
	"unicode"
)

// SnakeCase converts a string into snake case.
func SnakeCase(s string) string {
	return delimiterCase(s, '_', false)
}

// UpperSnakeCase converts a string into snake case with capital letters.
func UpperSnakeCase(s string) string {
	return delimiterCase(s, '_', true)
}

// delimiterCase converts a string into snake_case or kebab-case depending on the delimiter
func delimiterCase(s string, delimiter rune, upperCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	var builder strings.Builder
	builder.Grow(len(s)) // Pre-allocate space with some extra for delimiters

	adjustCase := unicode.ToLower
	if upperCase {
		adjustCase = unicode.ToUpper
	}

	runes := []rune(s)
	var prev rune

	for i := 0; i < len(runes); i++ {
		curr := runes[i]

		// Skip if current character is a delimiter
		if isDelimiterOptimized(curr) {
			if i > 0 && !isDelimiterOptimized(prev) {
				builder.WriteRune(delimiter)
			}
			prev = curr
			continue
		}

		// Handle consecutive uppercase letters (e.g., HTTP, UUID)
		if i > 0 && unicode.IsUpper(curr) {
			if unicode.IsLower(prev) {
				// Previous was lowercase, add delimiter (e.g., "fooBar" -> "foo_bar")
				builder.WriteRune(delimiter)
			} else if unicode.IsUpper(prev) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				// Previous was uppercase and next is lowercase (e.g., "FOOBar" -> "foo_bar")
				builder.WriteRune(delimiter)
			}
		}

		// Handle transitions between letters and numbers
		if unicode.IsNumber(curr) {
			if i > 0 && !unicode.IsNumber(prev) && !isDelimiterOptimized(prev) && unicode.IsUpper(prev) {
				// Only add delimiter when transitioning from an uppercase letter to number
				builder.WriteRune(delimiter)
			}
			builder.WriteRune(curr)
		} else {
			if i > 0 && unicode.IsNumber(prev) && unicode.IsUpper(curr) {
				// Add delimiter when transitioning from number to uppercase letter
				builder.WriteRune(delimiter)
			}
			builder.WriteRune(adjustCase(curr))
		}

		prev = curr
	}

	return builder.String()
}

// isDelimiterOptimized checks if a character is some kind of whitespace or '_' or '-'.
func isDelimiterOptimized(ch rune) bool {
	return ch == '-' || ch == '_' || ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}
