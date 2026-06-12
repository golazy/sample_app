package markdown

import (
	"strings"
	"testing"
)

func TestConvert(t *testing.T) {
	html, err := Convert("Hello, **world**.\n\n<script>unsafe()</script>")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(html, "<strong>world</strong>") {
		t.Fatalf("expected strong markup, got %q", html)
	}
	if strings.Contains(html, "<script>") {
		t.Fatalf("raw HTML was rendered: %q", html)
	}
}
