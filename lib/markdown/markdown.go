package markdown

import (
	"bytes"

	"github.com/yuin/goldmark"
)

var converter = goldmark.New()

// Convert renders Markdown as HTML. Goldmark's default configuration does not
// render raw HTML from the source.
func Convert(source string) (string, error) {
	var output bytes.Buffer
	if err := converter.Convert([]byte(source), &output); err != nil {
		return "", err
	}
	return output.String(), nil
}
