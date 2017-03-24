package kkok

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"time"

	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

var (
	// baseOtto is the source for all other Otto runtime.
	// baseOtto.Copy() will create a copy/clone of the runtime.
	baseOtto = otto.New()

	reFilterID = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// Filter is the interface that filter plugins must implement.
//
// Methods other than Params and Process are implemented in BaseFilter
// so that a filter implementation can embed BaseFilter to provide them.
type Filter interface {

	// Params returns PluginParams that can be used to construct
	// this filter.
	Params() PluginParams

	// Process applies the filter for all alerts and returns
	// the filtered alerts.
	Process(alerts []*Alert) ([]*Alert, error)

	// ID returns the ID of the filter.
	ID() string

	// Label returns the string label of the filter.
	Label() string

	// Dynamic returns true if the filter is added dynamically.
	Dynamic() bool

	// SetDynamic sets the filter dynamic.
	// This should be used privately inside kkok.
	SetDynamic()

	// Disabled returns true iff the filter is disabled.
	Disabled() bool

	// Enable enables the filter if e == true, otherwise disables the filter.
	Enable(e bool)

	// Expired returns true iff the filter has already been expired.
	Expired() bool
}

// FilterConstructor is a function signature for filter construction.
//
// id should be passed to BaseFilter.Init.
// params may be used to initialize the filter.
type FilterConstructor func(id string, params map[string]interface{}) (Filter, error)

var filterTypes = make(map[string]FilterConstructor)

// RegisterFilter registers a construction function of a Filter type.
func RegisterFilter(typ string, ctor FilterConstructor) {
	filterTypes[typ] = ctor
}

// NewFilter constructs a Filter.
func NewFilter(typ string, params map[string]interface{}) (Filter, error) {
	ctor, ok := filterTypes[typ]
	if !ok {
		return nil, errors.New("no such filter type: " + typ)
	}
	i, ok := params["id"]
	if !ok {
		return nil, errors.New("no filter id")
	}
	id, ok := i.(string)
	if !ok {
		return nil, errors.New("filter id must be a string")
	}
	delete(params, "id")
	return ctor(id, params)
}

// BaseFilter provides the common implementation for all filters.
//
// Filter plugins MUST embed BaseFilter anonymously.
type BaseFilter struct {
	id        string
	label     string
	dynamic   bool
	disabled  bool
	all       bool
	origIf    interface{}
	ifScript  *otto.Script
	ifCommand *exec.Cmd
	expire    time.Time
}

// Init initializes BaseFilter with parameters.
//
// Significant keys in params are:
//
//    label    string  Arbitrary string label of the filter.
//    disabled bool    If true, this filter is disabled.
//    expire   string  RFC3339 format time at which this filter expires.
//    all      bool    If true, the filter process all alerts at once.
//    if       string|array of strings
//                     string must be a JavaScript expression.
//                     array of strings must be a command and arguments
//                     to be invoked.
func (b *BaseFilter) Init(id string, params map[string]interface{}) error {
	if !reFilterID.MatchString(id) {
		return errors.New("invalid filter id: " + id)
	}

	b.id = id

	if i, ok := params["label"]; ok {
		label, ok := i.(string)
		if !ok {
			return errors.New("label must be a string")
		}
		b.label = label
	}

	if i, ok := params["disabled"]; ok {
		disabled, ok := i.(bool)
		if !ok {
			return errors.New("disabled must be a boolean")
		}
		b.disabled = disabled
	}

	if i, ok := params["expire"]; ok {
		s, ok := i.(string)
		if !ok {
			return errors.New("expire must be a string for RFC3339 time format")
		}
		err := b.expire.UnmarshalText([]byte(s))
		if err != nil {
			return errors.Wrap(err, "expire")
		}
	}

	if i, ok := params["all"]; ok {
		all, ok := i.(bool)
		if !ok {
			return errors.New("all must be a boolean")
		}
		b.all = all
	}

	if _, ok := params["if"]; !ok {
		return nil
	}

	b.origIf = params["if"]

	switch i := b.origIf.(type) {
	case string:
		script, err := baseOtto.Compile("", i)
		if err != nil {
			return errors.Wrap(err, "expr: "+i)
		}
		b.ifScript = script

	case []interface{}:
		command, err := parseCommand(i)
		if err != nil {
			return err
		}
		b.ifCommand = command

	default:
		return errors.New("if must be string or []string")
	}

	return nil
}

func parseCommand(a []interface{}) (*exec.Cmd, error) {
	if len(a) == 0 {
		return nil, errors.New("empty command")
	}

	b := make([]string, len(a))
	for i, e := range a {
		s, ok := e.(string)
		if !ok {
			return nil, fmt.Errorf("not a string: %#v", e)
		}
		b[i] = s
	}

	return exec.Command(b[0], b[1:]...), nil
}

// AddParams adds basic parameters for Filter.Params method.
// The map must not be nil.
func (b *BaseFilter) AddParams(m map[string]interface{}) {
	if b.disabled {
		m["disabled"] = b.disabled
	}
	if b.all {
		m["all"] = b.all
	}
	if b.origIf != nil {
		m["if"] = b.origIf
	}
	if b.dynamic && (!b.expire.IsZero()) {
		m["expire"] = b.expire
	}
}

// ID returns the ID of the filter.
func (b *BaseFilter) ID() string {
	return b.id
}

// Label returns the string label of the filter.
func (b *BaseFilter) Label() string {
	return b.label
}

// Dynamic returns true if the filter is added dynamically.
func (b *BaseFilter) Dynamic() bool {
	return b.dynamic
}

// SetDynamic sets the filter dynamic.
// This should be used privately inside kkok.
func (b *BaseFilter) SetDynamic() {
	b.dynamic = true
}

// Disabled returns true iff the filter is disabled.
func (b *BaseFilter) Disabled() bool {
	return b.disabled
}

// Enable enables the filter if e == true, otherwise disables the filter.
func (b *BaseFilter) Enable(e bool) {
	b.disabled = !e
}

// All returns true iff the filter processes all alerts at once.
// If false, the filter processes alerts one by one.
func (b *BaseFilter) All() bool {
	return b.all
}

// Expired returns true iff the filter has already been expired.
func (b *BaseFilter) Expired() bool {
	if !b.dynamic {
		return false
	}

	if b.expire.IsZero() {
		return false
	}

	return time.Now().After(b.expire)
}

// EvalAlert evaluates an alert with "if" condition.
func (b *BaseFilter) EvalAlert(a *Alert) (bool, error) {
	if b.ifScript != nil {
		vm := baseOtto.Copy()
		err := vm.Set("alert", a)
		if err != nil {
			return false, err
		}
		value, err := vm.Run(b.ifScript)
		if err != nil {
			return false, err
		}
		bvalue, _ := value.ToBoolean()
		return bvalue, nil
	}

	if b.ifCommand != nil {
		command := &exec.Cmd{
			Path: b.ifCommand.Path,
			Args: b.ifCommand.Args,
		}
		w, err := command.StdinPipe()
		if err != nil {
			return false, err
		}
		err = command.Start()
		if err != nil {
			return false, err
		}

		enc := json.NewEncoder(w)
		err = enc.Encode(a)
		w.Close()
		if err != nil {
			command.Wait()
			return false, err
		}

		return command.Wait() == nil, nil
	}

	return true, nil
}

// EvalAllAlerts evaluates all alerts with "if" condition.
func (b *BaseFilter) EvalAllAlerts(alerts []*Alert) (bool, error) {
	if b.ifScript != nil {
		vm := baseOtto.Copy()
		err := vm.Set("alerts", alerts)
		if err != nil {
			return false, err
		}
		value, err := vm.Run(b.ifScript)
		if err != nil {
			return false, err
		}
		bvalue, _ := value.ToBoolean()
		return bvalue, nil
	}

	if b.ifCommand != nil {
		command := &exec.Cmd{
			Path: b.ifCommand.Path,
			Args: b.ifCommand.Args,
		}
		return command.Run() == nil, nil
	}

	return true, nil
}
