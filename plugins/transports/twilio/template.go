package twilio

import "text/template"

// DefaultTemplate is the default text/template to render alert message body.
// "slack" is a template function to escape special characters in Slack.
const DefaultTemplate = `Title: {{.Title}}
From: {{.From}}
Host: {{.Host}}
Message: {{.Message}}`

var (
	defaultTemplate = template.Must(template.New("").Parse(DefaultTemplate))
)
