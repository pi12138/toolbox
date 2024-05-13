package head

import (
	"strings"
)

type Title struct {
	Content string
}

func (t Title) Render() string {
	var b strings.Builder
	b.WriteString("<title>")
	b.WriteString(t.Content)
	b.WriteString("</title>")
	return b.String()
}
