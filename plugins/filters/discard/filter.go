package discard

import "github.com/cybozu-go/kkok"

const filterType = "discard"

type filter struct {
	kkok.BaseFilter
}

func (f *filter) Params() kkok.PluginParams {
	m := make(map[string]interface{})
	f.BaseFilter.AddParams(m)
	return kkok.PluginParams{
		Type:   filterType,
		Params: m,
	}
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	if f.BaseFilter.All() {
		ok, err := f.BaseFilter.IfAll(alerts)
		if err != nil {
			return nil, err
		}

		if ok {
			return nil, nil
		}
		return alerts, nil
	}

	filtered := make([]*kkok.Alert, 0, len(alerts))
	for _, a := range alerts {
		ok, err := f.BaseFilter.If(a)
		if err != nil {
			return nil, err
		}

		if !ok {
			filtered = append(filtered, a)
		}
	}
	return filtered, nil
}

func ctor(id string, params map[string]interface{}) (kkok.Filter, error) {
	f := &filter{}
	err := f.Init(id, params)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func init() {
	kkok.RegisterFilter(filterType, ctor)
}
