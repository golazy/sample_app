package helpers

import (
	"strings"
	"testing"
)

func TestWordCount(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{name: "empty", content: "", want: 0},
		{name: "markdown", content: "Welcome to **GoLazy**, a controller-oriented Go app.", want: 7},
		{name: "unicode", content: "Buenos dias, manana cafe 2026", want: 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := WordCount(test.content); got != test.want {
				t.Fatalf("WordCount(%q) = %d, want %d", test.content, got, test.want)
			}
		})
	}
}

func TestReadTime(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    int
	}{
		{name: "empty", content: "", want: 0},
		{name: "one word", content: "GoLazy", want: 1},
		{name: "two hundred words", content: strings.Repeat("word ", 200), want: 1},
		{name: "rounds up", content: strings.Repeat("word ", 201), want: 2},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := ReadTime(test.content); got != test.want {
				t.Fatalf("ReadTime(%q) = %d, want %d", test.content, got, test.want)
			}
		})
	}
}

func TestRegisterHelpers(t *testing.T) {
	registered := RegisterHelpers()

	wordCount, ok := registered["word_count"].(func(string) int)
	if !ok {
		t.Fatalf("word_count helper = %T, want func(string) int", registered["word_count"])
	}
	if got := wordCount("one two"); got != 2 {
		t.Fatalf("word_count helper = %d, want 2", got)
	}

	readTime, ok := registered["read_time"].(func(string) int)
	if !ok {
		t.Fatalf("read_time helper = %T, want func(string) int", registered["read_time"])
	}
	if got := readTime("one two"); got != 1 {
		t.Fatalf("read_time helper = %d, want 1", got)
	}
}
