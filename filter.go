package kkok

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"sync"
	"time"

	"github.com/cybozu-go/kkok/util"
	"github.com/cybozu-go/log"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

var (
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

	// Disabled returns true if the filter is disabled or inactive.
	Disabled() bool

	// Enable enables the filter if e is true, otherwise disables the filter.
	// If the filter is currently inactive and e is true, then the filter
	// gets activated immediately.
	Enable(e bool)

	// Inactivate disables the filter until the given time.
	// Calling Enable(true) immediately re-enables the filter.
	Inactivate(until time.Time)

	// Expired returns true iff the filter has already been expired.
	Expired() bool

	// Reload reloads JavaScript files.
	Reload() error
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
	VM

	// constants
	id        string
	label     string
	dynamic   bool
	all       bool
	origIf    interface{}
	ifScript  *otto.Script
	ifCommand *exec.Cmd
	scripts   []string
	expire    time.Time

	// dynamic values
	mu            sync.Mutex
	disabled      bool
	inactiveUntil time.Time
}

// Init initializes BaseFilter with parameters.
//
// Significant keys in params are:
//
//    label    string    Arbitrary string label of the filter.
//    disabled bool      If true, this filter is disabled.
//    expire   string    RFC3339 format time at which this filter expires.
//    all      bool      If true, the filter process all alerts at once.
//    if       string | []string
//                       string must be a JavaScript expression.
//                       []string must be a command and arguments
//                       to be invoked.
//    scripts  []string  JavaScript filenames.
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

	scripts, err := util.GetStringSlice("scripts", params)
	switch {
	case err == nil:
		b.scripts = scripts
	case util.IsNotFound(err):
	default:
		return errors.Wrap(err, "scripts")
	}

	err = b.Reload()
	if err != nil {
		return errors.Wrap(err, "scripts")
	}

	if _, ok := params["if"]; !ok {
		return nil
	}

	b.origIf = params["if"]

	switch i := b.origIf.(type) {
	case string:
		script, err := CompileJS(i)
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
	if len(b.label) > 0 {
		m["label"] = b.label
	}
	if b.disabled {
		m["disabled"] = b.disabled
	}
	if b.all {
		m["all"] = b.all
	}
	if len(b.scripts) > 0 {
		m["scripts"] = b.scripts
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

// Disabled returns true if the filter is disabled or inactive.
func (b *BaseFilter) Disabled() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.disabled {
		return true
	}
	if b.inactiveUntil.IsZero() {
		return false
	}
	return time.Now().Before(b.inactiveUntil)
}

// Enable enables the filter if e is true, otherwise disables the filter.
// If the filter is currently inactive and e is true, then the filter
// gets activated immediately.
func (b *BaseFilter) Enable(e bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.disabled = !e
	if e {
		var zero time.Time
		b.inactiveUntil = zero
	}
}

// Inactivate disables the filter until the given time.
// Calling Enable(true) immediately re-enables the filter.
func (b *BaseFilter) Inactivate(until time.Time) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.inactiveUntil = until
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

// Reload reloads JavaScript files.
func (b *BaseFilter) Reload() error {
	if len(b.scripts) == 0 {
		b.VM = VM{baseOtto}
		return nil
	}

	b.VM = NewVM()
	return b.VM.Load(b.scripts)
}

// If evaluates an alert with "if" condition.
func (b *BaseFilter) If(a *Alert) (bool, error) {
	if b.ifScript != nil {
		value, err := b.VM.EvalAlert(a, b.ifScript)
		if err != nil {
			return false, err
		}
		if !value.IsBoolean() {
			log.Warn("kkok: not a boolean expression", map[string]interface{}{
				"id":         b.id,
				"expression": b.origIf,
			})
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

// IfAll evaluates all alerts with "if" condition.
func (b *BaseFilter) IfAll(alerts []*Alert) (bool, error) {
	if b.ifScript != nil {
		value, err := b.VM.EvalAlerts(alerts, b.ifScript)
		if err != nil {
			return false, err
		}
		if !value.IsBoolean() {
			log.Warn("kkok: not a boolean expression", map[string]interface{}{
				"id":         b.id,
				"expression": b.origIf,
			})
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
