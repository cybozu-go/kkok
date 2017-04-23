package edit

import (
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/pkg/errors"
	"github.com/robertkrimen/otto"
)

var convertScript *otto.Script

func toObject(a *kkok.Alert) (*otto.Object, error) {
	vm := kkok.NewVM()
	val, err := vm.EvalAlert(a, convertScript)
	if err != nil {
		return nil, err
	}
	return val.Object(), nil
}

func fromObject(obj *otto.Object) (*kkok.Alert, error) {
	i, err := obj.Value().Export()
	if err != nil {
		return nil, err
	}

	a := &kkok.Alert{}
	for k, v := range i.(map[string]interface{}) {
		switch k {
		case "From":
			s, ok := v.(string)
			if !ok {
				return nil, errors.New("From is not a string")
			}
			a.From = s
		case "Date":
			v, err := obj.Get("Date")
			if err != nil {
				return nil, err // unlikely
			}
			if !v.IsObject() {
				return nil, errors.New("Date is not a JavaScript object")
			}
			val, err := v.Object().Call("toISOString")
			if err != nil {
				return nil, errors.Wrap(err, "Date is not a JavaScript Date")
			}
			dt, err := time.Parse(time.RFC3339Nano, val.String())
			if err != nil {
				return nil, errors.Wrap(err, "time.Parse")
			}
			a.Date = dt
		case "Host":
			s, ok := v.(string)
			if !ok {
				return nil, errors.New("Host is not a string")
			}
			a.Host = s
		case "Title":
			s, ok := v.(string)
			if !ok {
				return nil, errors.New("Title is not a string")
			}
			a.Title = s
		case "Message":
			s, ok := v.(string)
			if !ok {
				return nil, errors.New("Message is not a string")
			}
			a.Message = s
		case "Routes":
			il, ok := v.([]interface{})
			if ok && len(il) == 0 {
				break
			}
			sl, ok := v.([]string)
			if !ok {
				return nil, errors.New("Routes is not []string")
			}
			a.Routes = sl
		case "Info":
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("Info is not map[string]interface{}")
			}
			for kk, vv := range m {
				if i64, ok := vv.(int64); ok {
					m[kk] = int(i64)
				}
			}
			a.Info = m
		case "Stats":
			m, ok := v.(map[string]interface{})
			if !ok {
				return nil, errors.New("Stats is not map[string]interface{}")
			}
			mf := make(map[string]float64, len(m))
			for kk, vv := range m {
				switch vv := vv.(type) {
				case float64:
					mf[kk] = vv
				case int64:
					mf[kk] = float64(vv)
				case int:
					mf[kk] = float64(vv)
				default:
					return nil, errors.New("non-integer value in Stats")
				}
			}
			a.Stats = mf
		case "Sub":
			sub, ok := v.([]*kkok.Alert)
			if ok {
				a.Sub = sub
			}
			// ignore otherwise
		}
	}

	err = a.Validate()
	if err != nil {
		return nil, err
	}

	return a, nil
}

func init() {
	s, err := kkok.CompileJS(`
routes = new Array();
for( i = 0; i < alert.Routes.length; i++) {
    routes.push(alert.Routes[i]);
}
info = {}
for (var k in alert.Info) {
    info[k] = alert.Info[k];
}
stats = {}
for (var k in alert.Stats) {
    stats[k] = alert.Stats[k]
}
sub = new Array();
for( i = 0; i < alert.Sub.length; i++) {
    sub.push(alert.Sub[i]);
}
({
    "From": alert.From,
    "Date": new Date(alert.Date.UTC().Format("2006-01-02T15:04:05.000Z07:00")),
    "Host": alert.Host,
    "Title": alert.Title,
    "Message": alert.Message,
    "Routes": routes,
    "Info": info,
    "Stats": stats,
    "Sub": sub,
})`)
	if err != nil {
		panic(err)
	}
	convertScript = s
}
