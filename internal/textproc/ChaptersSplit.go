package textproc

import "strings"

func SplitIntoChapters(text string) []string {
	chapters := strings.Builder{}
	result := []string{}
	strs := strings.Split(text, "\n")
	for _, str := range strs {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(str)), "глава") {
			if chapters.Len() > 0 {
				trimmed := strings.TrimSuffix(chapters.String(), "\n")
				result = append(result, trimmed)
				chapters.Reset()
			}
		}
		chapters.WriteString(str)
		chapters.WriteString("\n")
	}
	if chapters.Len() > 0 {
		trimmed := strings.TrimSuffix(chapters.String(), "\n")
		result = append(result, trimmed)
	}
	return result
}
