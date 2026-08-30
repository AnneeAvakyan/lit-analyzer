package charextraction

import (
	"strings"
	"unicode"
)

var stopWords = map[string]bool{
	"бог":       true,
	"господь":   true,
	"боже":      true,
	"господи":   true,
	"отечество": true,
	"родина":    true,
	"он":        true,
	"она":       true,
	"оно":       true,
	"они":       true,
}

func isCapitalized(word string) bool {
	if word == "" {
		return false
	}
	runes := []rune(word)
	return unicode.IsUpper(runes[0])
}

func isStopWord(word string) bool {
	return stopWords[strings.ToLower(word)]
}

// не учитывает падежи русского языка . Планируется объединение через
// ручное слияние алиасов (POST /characters/merge) на более позднем этапе.
func ExtractCandidatesWithFrequency(sentences []string) map[string]int {
	freq := make(map[string]int)

	for _, sentence := range sentences {
		candidates := extractCandidatesFromSentence(sentence)
		for _, candidate := range candidates {
			freq[candidate]++
		}
	}
	return freq
}

func FilterByFrequency(freq map[string]int, minCount int) map[string]int {
	result := make(map[string]int)
	for key, value := range freq {
		if value >= minCount {
			result[key] = value
		}
	}
	return result
}

func extractCandidatesFromSentence(sentence string) []string {
	candidates := []string{}
	currentSeries := []string{}

	flushSeries := func() {
		if len(currentSeries) > 0 {
			candidates = append(candidates, strings.Join(currentSeries, " "))
			currentSeries = []string{}
		}
	}

	for i, word := range strings.Fields(sentence) {
		if i == 0 {
			continue
		}

		cleaned := strings.TrimFunc(word, unicode.IsPunct)
		if cleaned == "" {
			continue
		}

		if isCapitalized(cleaned) && !isStopWord(cleaned) {
			currentSeries = append(currentSeries, cleaned)
		} else {
			flushSeries()
		}

	}
	flushSeries()

	return candidates
}
