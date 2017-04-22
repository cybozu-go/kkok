// Package all imports all filters to be compiled-in.
package all

import (
	// import all static plugins
	_ "github.com/cybozu-go/kkok/plugins/filters/discard"
	_ "github.com/cybozu-go/kkok/plugins/filters/edit"
	_ "github.com/cybozu-go/kkok/plugins/filters/exec"
	_ "github.com/cybozu-go/kkok/plugins/filters/freq"
	_ "github.com/cybozu-go/kkok/plugins/filters/group"
	_ "github.com/cybozu-go/kkok/plugins/filters/route"
)
