// Package all imports all static transport plugins.
package all

import (
	// import all static plugins
	_ "github.com/cybozu-go/kkok/plugins/transports/email"
	_ "github.com/cybozu-go/kkok/plugins/transports/slack"
)
