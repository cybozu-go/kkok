package group

import (
	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

func ctor(id string, params map[string]interface{}) (kkok.Filter, error) {
	f := &filter{}

	v, err := util.GetString("by", params)
	switch {
	case err == nil && len(v) > 0:
		s, err := kkok.CompileJS(v)
		if err != nil {
			return nil, errors.Wrap(err, "group: by")
		}
		f.by = s
		f.origBy = v
	case err == nil || util.IsNotFound(err):
	default:
		return nil, errors.Wrap(err, "group: by")
	}

	from, err := util.GetString("from", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "group: from")
	}
	f.from = from

	title, err := util.GetString("title", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "group: title")
	}
	f.title = title

	msg, err := util.GetString("message", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "group: message")
	}
	f.message = msg

	routes, err := util.GetStringSlice("routes", params)
	if err != nil && !util.IsNotFound(err) {
		return nil, errors.Wrap(err, "group: routes")
	}
	f.routes = routes

	err = f.Init(id, params)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	kkok.RegisterFilter(filterType, ctor)
}
