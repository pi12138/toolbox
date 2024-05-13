package head

import "strings"

type Meta struct {
	Charset string
	Content string
	Name    string
}

func (m Meta) Render() string {
	var b strings.Builder
	b.WriteString("<meta")
	if m.Charset != "" {
		b.WriteString(` charset="`)
		b.WriteString(m.Charset)
		b.WriteString(`"`)
	}
	if m.Content != "" {
		b.WriteString(` content="`)
		b.WriteString(m.Content)
		b.WriteString(`"`)
	}
	if m.Name != "" {
		b.WriteString(` name="`)
		b.WriteString(m.Name)
		b.WriteString(`"`)
	}
	b.WriteString(">")
	b.WriteString("</meta>")
	return b.String()
}
