package helpers

import "unicode"

const wordsPerMinute = 200

func RegisterHelpers() map[string]any {
	return map[string]any{
		"word_count": WordCount,
		"read_time":  ReadTime,
	}
}

func WordCount(content string) int {
	words := 0
	inWord := false
	for _, char := range content {
		if unicode.IsLetter(char) || unicode.IsDigit(char) {
			if !inWord {
				words++
				inWord = true
			}
			continue
		}
		if (char == '\'' || char == '-') && inWord {
			continue
		}
		inWord = false
	}
	return words
}

func ReadTime(content string) int {
	words := WordCount(content)
	if words == 0 {
		return 0
	}
	minutes := words / wordsPerMinute
	if words%wordsPerMinute != 0 {
		minutes++
	}
	return minutes
}
