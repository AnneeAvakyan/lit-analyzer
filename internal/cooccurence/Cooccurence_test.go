package cooccurence

import (
	"reflect"
	"testing"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

func TestBuildGraph(t *testing.T) {
	tests := []struct {
		name       string
		windowSize int
		input      []entities.Mention
		expected   EdgeWeight
	}{
		{
			name:       "simple",
			windowSize: 1,
			input: []entities.Mention{
				{CharacterID: 1, GlobalSentenceIndex: 0},
				{CharacterID: 2, GlobalSentenceIndex: 1},
			},
			expected: EdgeWeight{
				{1, 2}: 1,
			},
		},
		{
			name:       "clean case, 0 distance",
			windowSize: 1,
			input: []entities.Mention{
				{CharacterID: 1, GlobalSentenceIndex: 5},
				{CharacterID: 2, GlobalSentenceIndex: 5},
			},
			expected: EdgeWeight{
				{1, 2}: 1,
			},
		},
		{
			name:       "out of window",
			windowSize: 5,
			input: []entities.Mention{
				{CharacterID: 1, GlobalSentenceIndex: 0},
				{CharacterID: 2, GlobalSentenceIndex: 10},
			},
			expected: EdgeWeight{},
		},
		{
			name:       "more mentions",
			windowSize: 5,
			input: []entities.Mention{
				{CharacterID: 1, GlobalSentenceIndex: 0},
				{CharacterID: 2, GlobalSentenceIndex: 0},
				{CharacterID: 1, GlobalSentenceIndex: 3},
				{CharacterID: 2, GlobalSentenceIndex: 3},
			},
			expected: EdgeWeight{
				{1, 2}: 3,
			},
		},
		{
			name:       "mentions itself",
			windowSize: 5,
			input: []entities.Mention{
				{CharacterID: 1, GlobalSentenceIndex: 0},
				{CharacterID: 1, GlobalSentenceIndex: 1},
			},
			expected: EdgeWeight{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := BuildGraph(test.input, test.windowSize)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("expected: %v, got: %v", test.expected, result)
			}
		})
	}
}
