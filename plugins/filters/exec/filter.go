package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/cybozu-go/kkok"
	"github.com/cybozu-go/log"
	"github.com/cybozu-go/well"
	"github.com/pkg/errors"
)

const (
	filterType = "exec"

	defaultTimeout = 5 * time.Second
)

type filter struct {
	kkok.BaseFilter

	// constants
	command []string
	timeout time.Duration
}

func newFilter() *filter {
	return &filter{
		timeout: defaultTimeout,
	}
}

func (f *filter) Params() kkok.PluginParams {
	m := map[string]interface{}{
		"command": f.command,
		"timeout": int(f.timeout.Seconds()),
	}
	f.BaseFilter.AddParams(m)

	return kkok.PluginParams{
		Type:   filterType,
		Params: m,
	}
}

func (f *filter) exec(j []byte) ([]byte, error) {
	ctx := context.Background()
	if f.timeout != 0 {
		ctx2, cancel := context.WithTimeout(ctx, f.timeout)
		ctx = ctx2
		defer cancel()
	}

	command := well.CommandContext(ctx, f.command[0], f.command[1:]...)
	// suppress successful log
	command.Severity = log.LvDebug
	command.Stdin = bytes.NewReader(j)
	return command.Output()
}

func (f *filter) Process(alerts []*kkok.Alert) ([]*kkok.Alert, error) {
	var newAlerts []*kkok.Alert

	if f.BaseFilter.All() {
		ok, err := f.BaseFilter.IfAll(alerts)
		if err != nil {
			return nil, errors.Wrap(err, "exec:"+f.ID())
		}

		if !ok {
			return alerts, nil
		}

		j, err := json.Marshal(alerts)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal(alerts)")
		}

		jj, err := f.exec(j)
		if err != nil {
			return nil, errors.Wrap(err, "exec:"+f.ID())
		}

		err = json.Unmarshal(jj, &newAlerts)
		if err != nil {
			return nil, errors.Wrap(err, "exec:"+f.ID())
		}
		return newAlerts, nil
	}

	for _, a := range alerts {
		ok, err := f.BaseFilter.If(a)
		if err != nil {
			return nil, errors.Wrap(err, "exec:"+f.ID())
		}

		if !ok {
			newAlerts = append(newAlerts, a)
			continue
		}

		j, err := json.Marshal(a)
		if err != nil {
			return nil, errors.Wrap(err, "json.Marshal(a)")
		}

		jj, err := f.exec(j)
		if err != nil {
			return nil, errors.Wrap(err, "exec:"+f.ID())
		}

		aa := new(kkok.Alert)
		err = json.Unmarshal(jj, aa)
		if err != nil {
			return nil, errors.Wrap(err, "exec:"+f.ID())
		}
		newAlerts = append(newAlerts, aa)
	}

	return newAlerts, nil
}
