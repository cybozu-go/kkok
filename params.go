package kkok

import (
	"encoding/json"
	"fmt"

	"github.com/cybozu-go/kkok/util"
	"github.com/pkg/errors"
)

// PluginParams is used to construct plugins including filters and transports.
type PluginParams struct {
	Type   string
	Params map[string]interface{}
}

// UnmarshalTOML is to load PluginParams from TOML file.
func (t *PluginParams) UnmarshalTOML(i interface{}) error {
	data, ok := i.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid type: %T", i)
	}
	tt, err := util.GetString("type", data)
	if err != nil {
		return errors.Wrap(err, "failed to get type")
	}
	delete(data, "type")

	t.Type = tt
	t.Params = data
	return nil
}

// MarshalJSON implements json.Marshaler.
func (t PluginParams) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{}, len(t.Params)+1)
	for k, v := range t.Params {
		m[k] = v
	}
	m["type"] = t.Type
	return json.Marshal(m)
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *PluginParams) UnmarshalJSON(data []byte) error {
	var m map[string]interface{}
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}

	tt, err := util.GetString("type", m)
	if err != nil {
		return errors.Wrap(err, "failed to get type")
	}
	delete(m, "type")

	t.Type = tt
	t.Params = m
	return nil
}
