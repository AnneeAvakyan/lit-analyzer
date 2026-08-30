package textproc

import (
	"strings"
	"unicode"
)

var abbreviations = map[string]bool{
	"т.д.": true,
	"т.е.": true,
	"т.к.": true,
	"т.п.": true,
	"г.":   true,
	"гг.":  true,
	"см.":  true,
	"др.":  true,
}

var endings = map[rune]bool{
	'.': true,
	'?': true,
	'!': true,
}

func Segment(text string) []string {
	runes := []rune(text)
	result := []string{}
	if len(runes) == 0 {
		return []string{}
	}

	buffer := strings.Builder{}

	for i := 0; i < len(runes); i++ {
		buffer.WriteRune(runes[i])
		if !endings[runes[i]] {
			continue
		}
		if isAbbreviation(buffer.String()) {
			continue
		}
		peekWord := strings.Fields(buffer.String())
		if len(peekWord) > 0 {
			lastWord := strings.TrimSuffix(peekWord[len(peekWord)-1], ".")
			if isInitial(lastWord) {
				continue
			}
		}

		next := peekNextNonSpace(runes, i+1)
		if next != 0 && unicode.IsLower(next) {
			continue
		}

		result = append(result, strings.TrimSpace(buffer.String()))
		buffer.Reset()
	}
	if buffer.Len() > 0 {
		result = append(result, strings.TrimSpace(buffer.String()))
	}
	return result
}

func isAbbreviation(buffer string) bool {
	words := strings.Fields(buffer)
	if len(words) == 0 {
		return false
	}

	lastWord := words[len(words)-1]

	if abbreviations[strings.ToLower(lastWord)] {
		return true
	}

	return false
}

func isInitial(word string) bool {
	runes := []rune(word)
	if len(runes) != 1 {
		return false
	}

	return unicode.IsUpper(runes[0])
}

func peekNextNonSpace(runes []rune, index int) rune {
	for i := index; i < len(runes); i++ {
		if unicode.IsSpace(runes[i]) {
			continue
		}
		return runes[i]
	}
	return 0
}
