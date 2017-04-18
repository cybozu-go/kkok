package slack

import (
	"strings"
	"text/template"
)

// DefaultTemplate is the default text/template to render alert message body.
// "slack" is a template function to escape special characters in Slack.
const DefaultTemplate = `{{slack .Message}}`

var (
	// EscapeSlack is a string replacer for special characters in Slack.
	EscapeSlack = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace

	defaultTemplate = template.Must(newTemplate().Parse(DefaultTemplate))
)

func newTemplate() *template.Template {
	return template.New("").Funcs(map[string]interface{}{
		"slack": EscapeSlack,
	})
}
