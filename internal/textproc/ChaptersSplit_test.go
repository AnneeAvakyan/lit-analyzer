package textproc

import (
	"reflect"
	"testing"
)

func TestSplitIntoChapters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "two chapters",
			input: "Глава 1\nЛето кончилось.\nГлава 2.\nНенавижу бауманку.",
			expected: []string{
				"Глава 1\nЛето кончилось.",
				"Глава 2.\nНенавижу бауманку.",
			},
		},
		{
			name:     "no chapters",
			input:    "Лето кончилось.\nНенавижу бауманку.",
			expected: []string{"Лето кончилось.\nНенавижу бауманку."},
		},
		{
			name:  "two chapters with capslock",
			input: "ГЛАВА 1\nЛето кончилось.\nГлава 2.\nНенавижу бауманку.",
			expected: []string{
				"ГЛАВА 1\nЛето кончилось.",
				"Глава 2.\nНенавижу бауманку.",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := SplitIntoChapters(test.input)
			if !reflect.DeepEqual(test.expected, result) {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
