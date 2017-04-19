package kkok

import (
	"os"
	"path/filepath"

	"github.com/robertkrimen/otto"
)

var (
	// baseOtto is the source for all other Otto runtime.
	// baseOtto.Copy() will create a copy/clone of the runtime.
	baseOtto = otto.New()

	undefined = otto.UndefinedValue()
)

// CompileJS compiles a JavaScript expression s.
// The returned script can be passed to VM.EvalAlert and VM.EvalAlerts.
func CompileJS(s string) (*otto.Script, error) {
	return baseOtto.Compile("", s)
}

// VM wraps otto JavaScript engine to provide convenient methods.
type VM struct {
	*otto.Otto
}

// NewVM creates a new JavaScript virtual machine.
func NewVM() VM {
	return VM{baseOtto.Copy()}
}

// Load reads JavaScript files and executes them.
// This alters VM states for later use.
func (vm VM) Load(filenames []string) error {
	for _, fn := range filenames {
		f, err := os.Open(fn)
		if err != nil {
			return err
		}
		s, err := vm.Otto.Compile(filepath.Base(fn), f)
		f.Close()
		if err != nil {
			return err
		}
		_, err = vm.Otto.Run(s)
		if err != nil {
			return err
		}
	}
	return nil
}

// EvalAlert evaluates a JavaScript script.
// Before evaluation, a is assigned to "alert" variable.
// This will not alter any VM state.
func (vm VM) EvalAlert(a *Alert, s *otto.Script) (otto.Value, error) {
	o := vm.Otto.Copy()
	err := o.Set("alert", a)
	if err != nil {
		return undefined, err
	}
	return o.Run(s)
}

// EvalAlerts evaluates a JavaScript script.
// Before evaluation, alerts are assigned to "alerts" variable.
// This will not alter any VM state.
func (vm VM) EvalAlerts(alerts []*Alert, s *otto.Script) (otto.Value, error) {
	o := vm.Otto.Copy()
	err := o.Set("alerts", alerts)
	if err != nil {
		return undefined, err
	}
	return o.Run(s)
}
