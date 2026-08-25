package charextraction

import (
	"reflect"
	"testing"
)

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
