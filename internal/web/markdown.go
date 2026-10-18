package web

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
)

var newsMarkdown = goldmark.New()

// renderNewsMarkdown uses goldmark's safe defaults: raw HTML and dangerous
// link schemes are not rendered. The resulting HTML may therefore be passed
// to html/template as trusted renderer output.
func renderNewsMarkdown(source string) (template.HTML, error) {
	var output bytes.Buffer
	if err := newsMarkdown.Convert([]byte(source), &output); err != nil {
		return "", err
	}
	return template.HTML(output.String()), nil
}
