package exec

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/cybozu-go/cmd"
	"github.com/cybozu-go/kkok"
	"github.com/pkg/errors"
)

const (
	transportType = "exec"

	defaultTimeout = 5 * time.Second
)

type transport struct {
	label   string
	command string
	args    []string
	all     bool
	timeout time.Duration
}

func newTransport(command string, args ...string) *transport {
	return &transport{
		command: command,
		args:    args,
		timeout: defaultTimeout,
	}
}

func (t *transport) String() string {
	if len(t.label) > 0 {
		return t.label
	}

	return transportType
}

func (t *transport) Params() kkok.PluginParams {
	ca := make([]string, len(t.args)+1)
	ca[0] = t.command
	copy(ca[1:], t.args)
	m := map[string]interface{}{
		"command": ca,
		"timeout": int(t.timeout.Seconds()),
	}

	if len(t.label) > 0 {
		m["label"] = t.label
	}
	if t.all {
		m["all"] = t.all
	}

	return kkok.PluginParams{
		Type:   transportType,
		Params: m,
	}
}

func (t *transport) exec(j []byte) error {
	ctx := context.Background()
	if t.timeout != 0 {
		ctx2, cancel := context.WithTimeout(ctx, t.timeout)
		ctx = ctx2
		defer cancel()
	}

	command := cmd.CommandContext(ctx, t.command, t.args...)
	command.Stdin = bytes.NewReader(j)
	return command.Run()
}

func (t *transport) Deliver(alerts []*kkok.Alert) error {
	if t.all {
		data, err := json.Marshal(alerts)
		if err != nil {
			return errors.Wrap(err, transportType)
		}
		err = t.exec(data)
		if err != nil {
			return errors.Wrap(err, transportType)
		}
		return nil
	}

	for _, a := range alerts {
		data, err := json.Marshal(a)
		if err != nil {
			return errors.Wrap(err, transportType)
		}
		err = t.exec(data)
		if err != nil {
			return errors.Wrap(err, transportType)
		}
	}

	return nil
}
