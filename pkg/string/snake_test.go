package snake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSnakeCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"simple_string", "simple_string"},
		{"CamelCase", "camel_case"},
		{"ThisIsATest", "this_is_a_test"},
		{"HTTPRequest", "http_request"},
		{"UpperCamelCase", "upper_camel_case"},
		{"  leadingAndTrailingSpaces  ", "leading_and_trailing_spaces"},
		{"MixedCaseAnd_underscores", "mixed_case_and_underscores"},
		{"ALL_CAPS_STRING", "all_caps_string"},
		{"already_snake_case", "already_snake_case"},
		{"123NumbersAtStart", "123_numbers_at_start"},
		{"NumbersInMiddle123Word", "numbers_in_middle123_word"},
		{"WordWithNumbersAtEnd123", "word_with_numbers_at_end123"},
		{"SomeID", "some_id"},
		{"SomeUUID", "some_uuid"},
		{"unicodeTest你好世界", "unicode_test你好世界"}, // Unicode test
		{"  你好世界  ", "你好世界"},                    // Unicode with spaces
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := SnakeCase(tc.input)
			assert.Equal(t, tc.expected, actual, "SnakeCase(%q)", tc.input)
		})
	}
}

func TestUpperSnakeCase(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"simple_string", "SIMPLE_STRING"},
		{"CamelCase", "CAMEL_CASE"},
		{"ThisIsATest", "THIS_IS_A_TEST"},
		{"HTTPRequest", "HTTP_REQUEST"},
		{"UpperCamelCase", "UPPER_CAMEL_CASE"},
		{"  leadingAndTrailingSpaces  ", "LEADING_AND_TRAILING_SPACES"},
		{"MixedCaseAnd_underscores", "MIXED_CASE_AND_UNDERSCORES"},
		{"ALL_CAPS_STRING", "ALL_CAPS_STRING"},
		{"already_snake_case", "ALREADY_SNAKE_CASE"},
		{"123NumbersAtStart", "123_NUMBERS_AT_START"},
		{"NumbersInMiddle123Word", "NUMBERS_IN_MIDDLE123_WORD"},
		{"WordWithNumbersAtEnd123", "WORD_WITH_NUMBERS_AT_END123"},
		{"SomeID", "SOME_ID"},
		{"SomeUUID", "SOME_UUID"},
		{"unicodeTest你好世界", "UNICODE_TEST你好世界"}, // Unicode test
		{"  你好世界  ", "你好世界"},                    // Unicode with spaces
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			actual := UpperSnakeCase(tc.input)
			assert.Equal(t, tc.expected, actual, "UpperSnakeCase(%q)", tc.input)
		})
	}
}

func FuzzSnakeCase(f *testing.F) {
	testCases := []string{
		"",
		"simple_string",
		"CamelCase",
		"ThisIsATest",
		"HTTPRequest",
		"UpperCamelCase",
		"  leadingAndTrailingSpaces  ",
		"MixedCaseAnd_underscores",
		"ALL_CAPS_STRING",
		"already_snake_case",
		"123NumbersAtStart",
		"NumbersInMiddle123Word",
		"WordWithNumbersAtEnd123",
		"SomeID",
		"SomeUUID",
		"unicodeTest你好世界",
		"  你好世界  ",
	}
	for _, tc := range testCases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, input string) {
		SnakeCase(input) // Just run it, check for panics or errors
	})
}

func FuzzUpperSnakeCase(f *testing.F) {
	testCases := []string{
		"",
		"simple_string",
		"CamelCase",
		"ThisIsATest",
		"HTTPRequest",
		"UpperCamelCase",
		"  leadingAndTrailingSpaces  ",
		"MixedCaseAnd_underscores",
		"ALL_CAPS_STRING",
		"already_snake_case",
		"123NumbersAtStart",
		"NumbersInMiddle123Word",
		"WordWithNumbersAtEnd123",
		"SomeID",
		"SomeUUID",
		"unicodeTest你好世界",
		"  你好世界  ",
	}
	for _, tc := range testCases {
		f.Add(tc)
	}
	f.Fuzz(func(t *testing.T, input string) {
		UpperSnakeCase(input) // Just run it, check for panics or errors
	})
}
