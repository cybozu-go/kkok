package route

import (
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

const (
	defaultMuteSeconds = 60
)

func ctor(id string, params map[string]interface{}) (kkok.Filter, error) {
	routes, err := util.GetStringSlice("routes", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "route: routes")
	}

	muteRoutes, err := util.GetStringSlice("mute_routes", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "route: mute_routes")
	}

	replace, err := util.GetBool("replace", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "route: replace")
	}

	autoMute, err := util.GetBool("auto_mute", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "route: auto_mute")
	}

	muteSeconds, err := util.GetInt("mute_seconds", params)
	switch {
	case err == nil:
		if muteSeconds <= 0 {
			return nil, errors.New("route: invalid mute_seconds")
		}
	case util.IsNotFound(err):
		muteSeconds = defaultMuteSeconds
	default:
		return nil, errors.Wrap(err, "route: mute_seconds")
	}

	f := &filter{
		routes:       routes,
		replace:      replace,
		autoMute:     autoMute,
		muteDuration: time.Duration(muteSeconds) * time.Second,
		muteRoutes:   muteRoutes,
	}
	err = f.Init(id, params)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	kkok.RegisterFilter("route", ctor)
}
