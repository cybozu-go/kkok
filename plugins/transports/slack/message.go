package slack

import (
	"fmt"
	"strings"
)

const (
	maxShortLength = 24
)

// message is to encode JSON for Slack incoming message.
type message struct {
	Text        string        `json:"text,omitempty"`
	Name        string        `json:"username,omitempty"`
	Icon        string        `json:"icon_emoji,omitempty"`
	Channel     string        `json:"channel,omitempty"`
	Attachments []*attachment `json:"attachments,omitempty"`
}

// attachment is to encode JSON for Slack attachments.
type attachment struct {
	Fallback string   `json:"fallback"`
	Color    string   `json:"color,omitempty"`
	Title    string   `json:"title,omitempty"`
	Fields   []*field `json:"fields,omitempty"`
	Text     string   `json:"text,omitempty"`
}

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// addField adds a text field.  value will be automatically
// converted to a string if it is not.  Note that strings will
// be escaped by EscapeSlack automatically.
func (a *attachment) addField(title string, value interface{}) {
	f := &field{Title: title}
	switch v := value.(type) {
	case nil:
		return
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		f.Value = fmt.Sprint(v)
		f.Short = true
	case string:
		f.Value = EscapeSlack(v)
		f.Short = (len(v) <= maxShortLength) && (!strings.Contains(v, "\n"))
	default:
		f.Value = EscapeSlack(fmt.Sprintf("%#v", v))
		f.Short = false
	}
	a.Fields = append(a.Fields, f)
}
