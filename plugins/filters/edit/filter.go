package edit

import (
	"github.com/cybozu-go/kkok"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

const (
	filterType = "edit"
)

type filter struct {
	kkok.BaseFilter

	// constants
	code     *otto.Script
	origCode string
}

func (f *filter) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"code": f.origCode,
	}
	f.BaseFilter.AddParams(m)

	return kkok.PluginParams{
		Type:   filterType,
		Params: m,
	}
}

func (f *filter) edit(a *kkok.Alert) (*kkok.Alert, error) {
	obj, err := toObject(a)
	if err != nil {
		return nil, err
	}

	vm := kkok.NewVM()
	vm.Set("alert", obj.Value())
	_, err = vm.Run(f.code)
	if err != nil {
		return nil, err
	}

	return fromObject(obj)
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	for i, a := range alerts {
		ok, err := f.BaseFilter.If(a)
		if err != nil {
			return nil, errors.Wrap(err, "edit:"+f.ID())
		}

		if !ok {
			continue
		}

		aa, err := f.edit(a)
		if err != nil {
			return nil, errors.Wrap(err, "edit:"+f.ID())
		}
		alerts[i] = aa
	}

	return alerts, nil
}
