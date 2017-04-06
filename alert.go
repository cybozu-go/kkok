package kkok

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
)

const (
	maxFromLength  = 100
	maxTitleLength = 250
)

var (
	undefined = otto.UndefinedValue()
)

// Alert represents an alert.
type Alert struct {

	// From is an identifying string who sent this alert.
	// Example: "NTP monitor"
	From string

	// Date is the time when this alert is generated.
	Date time.Time

	// Host is the hostname or IP address where this alert is generated.
	Host string

	// Title is one-line description of the alert.
	Title string

	// Message is multi-line description of the alert.
	Message string `json:",omitempty"`

	// Routes contain route ID strings along which this alert is delivered.
	Routes []string

	// Info is a map of additional alert properties.
	Info map[string]interface{} `json:",omitempty"`

	// Stats is a map of dynamically calculated values by filters.
	// This field is ignored for JSON.
	Stats map[string]float64 `json:"-"`

	// Sub may list alerts grouped into this.
	Sub []*Alert `json:",omitempty"`
}

// Validate validates constructed Alert struct.
// For invalid structs, non-nil errors are returned.
func (a *Alert) Validate() error {
	if len(a.From) == 0 {
		return errors.New("empty From")
	}
	if len(a.From) > maxFromLength {
		return errors.New("too long From")
	}
	if strings.Contains(a.From, "\n") {
		return errors.New("multi-line From")
	}

	if len(a.Title) == 0 {
		return errors.New("empty Title")
	}
	if len(a.Title) > maxTitleLength {
		return errors.New("too long Title")
	}
	if strings.Contains(a.Title, "\n") {
		return errors.New("multi-line From")
	}

	return nil
}

// SetInfo sets a value in Info with key.
func (a *Alert) SetInfo(key string, value interface{}) {
	if a.Info == nil {
		a.Info = make(map[string]interface{})
	}
	a.Info[key] = value
}

// SetStat sets a statistics value with key.
func (a *Alert) SetStat(key string, value float64) {
	if a.Stats == nil {
		a.Stats = make(map[string]float64)
	}
	a.Stats[key] = value
}

// Clone returns a deeply-copied clone of a.
//
// Stats field is not copied.
func (a *Alert) Clone() *Alert {
	var croutes []string
	if len(a.Routes) > 0 {
		croutes = make([]string, len(a.Routes))
		copy(croutes, a.Routes)
	}

	var cinfo map[string]interface{}
	if len(a.Info) > 0 {
		cinfo = make(map[string]interface{})
		for k, v := range a.Info {
			cinfo[k] = v
		}
	}

	var csub []*Alert
	if len(a.Sub) > 0 {
		csub = make([]*Alert, len(a.Sub))
		for i, a2 := range a.Sub {
			csub[i] = a2.Clone()
		}
	}

	return &Alert{
		From:    a.From,
		Date:    a.Date,
		Host:    a.Host,
		Title:   a.Title,
		Message: a.Message,
		Routes:  croutes,
		Info:    cinfo,
		Sub:     csub,
	}
}

// String returns a string representation of the alert.
func (a *Alert) String() string {
	return fmt.Sprintf("[%s@%s] %s", a.From, a.Host, a.Title)
}

// Eval evaluates a JavaScript script.
// Alert itself is set to "alert" variable in the expression.
func (a *Alert) Eval(scr *otto.Script) (otto.Value, error) {
	vm := baseOtto.Copy()
	err := vm.Set("alert", a)
	if err != nil {
		return undefined, err
	}
	return vm.Run(scr)
}
