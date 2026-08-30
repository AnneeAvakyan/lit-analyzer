package textproc

import (
	"reflect"
	"testing"
)

func TestIsAbbreviation(t *testing.T) {
	tests := []struct {
		name     string
		buffer   string
		expected bool
	}{
		{
			name:     "known abbreviation lowercase",
			buffer:   "она поехала в москву и т.д.",
			expected: true,
		},
		{
			name:     "not an abbreviation",
			buffer:   "она поехала в москву",
			expected: false,
		},
		{
			name:     "unknown abbreviation",
			buffer:   "крч. дело было так",
			expected: false,
		},
		{
			name:     "empty buffer",
			buffer:   "",
			expected: false,
		},
		{
			name:     "uppercase",
			buffer:   "она поехала в москву Т.Д.",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAbbreviation(tt.buffer)
			if result != tt.expected {
				t.Errorf("isAbbreviation(%q) = %v, want %v", tt.buffer, result, tt.expected)
			}
		})
	}
}

func TestIsInitial(t *testing.T) {
	tests := []struct {
		name     string
		buffer   string
		expected bool
	}{
		{
			name:     "is initial",
			buffer:   "А",
			expected: true,
		},
		{
			name:     "not initial",
			buffer:   "Она",
			expected: false,
		},
		{
			name:     "empty buffer",
			buffer:   "",
			expected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isInitial(tt.buffer)
			if result != tt.expected {
				t.Errorf("isInitial(%q) = %v, want %v", tt.buffer, result, tt.expected)
			}
		})
	}
}

func TestPeekNextNonSpace(t *testing.T) {
	tests := []struct {
		name     string
		buffer   []rune
		index    int
		expected rune
	}{
		{
			name:     "empty buffer",
			buffer:   []rune(""),
			index:    0,
			expected: 0,
		},
		{
			name:     "full buffer with space",
			buffer:   []rune("Вечерело. Небо окрасилось в сирень."),
			index:    9,
			expected: 'Н',
		},
		{
			name:     "full buffer without space",
			buffer:   []rune("Август."),
			index:    7,
			expected: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := peekNextNonSpace(tt.buffer, tt.index)
			if result != tt.expected {
				t.Errorf("peekNextNonSpace(%q) = %v, want %v", tt.buffer, result, tt.expected)
			}
		})
	}
}

func TestSegment(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		expects []string
	}{
		{
			name:    "empty text",
			text:    "",
			expects: []string{},
		},
		{
			name:    "full text",
			text:    "Он был панком, а она балериной.",
			expects: []string{"Он был панком, а она балериной."},
		},
		{
			name:    "two sentences",
			text:    "Он был панком. Она - балериной.",
			expects: []string{"Он был панком.", "Она - балериной."},
		},
		{
			name:    "with initial",
			text:    "А. Загадочный был панком, а она балериной.",
			expects: []string{"А. Загадочный был панком, а она балериной."},
		},
		{
			name:    "with abbreviation",
			text:    "Он был панком и т.д., а она балериной.",
			expects: []string{"Он был панком и т.д., а она балериной."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Segment(tt.text)
			if !reflect.DeepEqual(result, tt.expects) {
				t.Errorf("Segment(%q) = %v, want %v", tt.text, result, tt.expects)
			}
		})
	}
}
