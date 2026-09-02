package charextraction

import (
	"reflect"
	"testing"
)

func TestExtractCandidatesOccurrences(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected map[string][]int
	}{
		{
			name:  "two sentences, freq = one & two",
			input: []string{"Юная Наташа полюбила Андрея.", "Князя Андрея это мало заботило."},
			expected: map[string][]int{
				"Наташа": {0},
				"Андрея": {0, 1},
			},
		},
		{
			name:  "two sentences, freq = one & one",
			input: []string{"Юная Наташа полюбила Андрея.", "Он отвечал ей взаимностью."},
			expected: map[string][]int{
				"Наташа": {0},
				"Андрея": {0},
			},
		},
		{
			name:     "two sentences, no names",
			input:    []string{"Она полюбила его.", "Полюбила на всю жизнь."},
			expected: map[string][]int{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := ExtractCandidatesOccurrences(test.input)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Expected: %v, Actual: %v", test.expected, actual)
			}
		})
	}
}

func TestFilterByFrequency(t *testing.T) {
	tests := []struct {
		name     string
		freq     map[string][]int
		minCount int
		expected map[string][]int
	}{
		{
			name: "remove less than 2",
			freq: map[string][]int{
				"Наташа": {1, 2},
				"Андрей": {1},
			},
			minCount: 2,
			expected: map[string][]int{
				"Наташа": {1, 2},
			},
		},
		{
			name: "no remove",
			freq: map[string][]int{
				"Наташа": {2},
				"Андрей": {1},
			},
			minCount: 1,
			expected: map[string][]int{
				"Наташа": {2},
				"Андрей": {1},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := FilterByFrequency(test.freq, test.minCount)
			if !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("Expected: %v, Actual: %v", test.expected, actual)
			}
		})
	}
}

func TestExtractCandidatesFromSentence(t *testing.T) {
	tests := []struct {
		name     string
		sentence string
		want     []string
	}{
		{
			name:     "simple complex name",
			sentence: "Она встретила Андрея Болконского.",
			want:     []string{"Андрея Болконского"},
		},
		{
			name:     "two complex names",
			sentence: "Она встретила Андрея Болконского и Наташу Ростову.",
			want:     []string{"Андрея Болконского", "Наташу Ростову"},
		},
		{
			name:     "skip name",
			sentence: "Наташа пошла домой.",
			want:     []string{},
		},
		{
			name:     "no names",
			sentence: "Она пошла домой.",
			want:     []string{},
		},
		{
			name:     "stop word end series",
			sentence: "Она Андрея Бог знает где искала.",
			want:     []string{"Андрея"},
		},
		{
			name:     "sign after name",
			sentence: "Она встретила Андрея, который направлялся домой.",
			want:     []string{"Андрея"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractCandidatesFromSentence(tt.sentence)
			if !reflect.DeepEqual(result, tt.want) {
				t.Errorf("extractCandidatesFromSentence() = %v, want %v", result, tt.want)
			}
		})
	}
}
