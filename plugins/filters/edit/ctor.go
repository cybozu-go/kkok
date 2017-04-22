package edit

import (
	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

func ctor(id string, params map[string]interface{}) (kkok.Filter, error) {
	s, err := util.GetString("code", params)
	if err != nil {
		return nil, errors.Wrap(err, "edit:"+id)
	}

	code, err := kkok.CompileJS(s)
	if err != nil {
		return nil, errors.Wrap(err, "edit:"+id)
	}

	all, _ := util.GetBool("all", params)
	if all {
		return nil, errors.New("edit:" + id + ": all is not supported")
	}

	f := &filter{
		code:     code,
		origCode: s,
	}
	err = f.Init(id, params)
	if err != nil {
		return nil, errors.Wrap(err, "edit:"+id)
	}

	return f, nil
}

func init() {
	kkok.RegisterFilter(filterType, ctor)
}
