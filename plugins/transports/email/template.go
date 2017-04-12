package email

import "text/template"

var (
	defaultTemplate = template.Must(template.New("").Parse(DefaultTemplate))
)

// DefaultTemplate is the default text/template template for mail body.
const DefaultTemplate = `{{block "body" . -}}
From: {{.From}}
Date: {{.Date.UTC.Format "2006-01-02T15:04:05.999999999Z07:00"}}
Host: {{.Host}}
Title: {{.Title}}
{{if .Message}}
{{.Message -}}
{{end -}}
{{if .Info}}
Info:
{{range $key, $value := .Info}}  {{$key}}={{$value}}
{{end}}{{end -}}
{{if .Sub}}
Sub alerts:
-------------------------------------------------------
{{range .Sub}}{{template "body" . -}}
-------------------------------------------------------
{{end}}{{end}}{{end}}`
