package kkok

import "github.com/robertkrimen/otto"

var (
	// baseOtto is the source for all other Otto runtime.
	// baseOtto.Copy() will create a copy/clone of the runtime.
	baseOtto = otto.New()
)

// CompileJS compiles a JavaScript expression s.
// The returned script can be passed to Alert.Eval later.
func CompileJS(s string) (*otto.Script, error) {
	return baseOtto.Compile("", s)
}
