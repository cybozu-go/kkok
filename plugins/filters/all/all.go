// Package all imports all filters to be compiled-in.
package all

import (
	// import all static plugins
	_ "github.com/cybozu-go/kkok/plugins/filters/freq"
	_ "github.com/cybozu-go/kkok/plugins/filters/route"
)
