package kkok

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

// PluginParams is used to construct plugins including filters and transports.
type PluginParams struct {
	Type   string
	Params map[string]interface{}
}

func getString(k string, m map[string]interface{}) (string, error) {
	v, ok := m[k]
	if !ok {
		return "", errors.New("no such key: " + k)
	}
	s, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("value for %s is not string", k)
	}
	return s, nil
}

// UnmarshalTOML is to load PluginParams from TOML file.
func (t *PluginParams) UnmarshalTOML(i interface{}) error {
	data, ok := i.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid type: %T", i)
	}
	tt, err := getString("type", data)
	if err != nil {
		return err
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

	tt, err := getString("type", m)
	if err != nil {
		return err
	}
	delete(m, "type")

	t.Type = tt
	t.Params = m
	return nil
}

// GetAsString returns a string parameter associated with k in t.Params.
func (t *PluginParams) GetAsString(k string) (string, error) {
	i, ok := t.Params[k]
	if !ok {
		return "", errors.New("no such parameter: " + k)
	}
	s, ok := i.(string)
	if !ok {
		return "", fmt.Errorf("%s is not a string: %T", k, i)
	}
	return s, nil
}

// GetAsFloat returns a number parameter associated with k in t.Params.
func (t *PluginParams) GetAsFloat(k string) (float64, error) {
	i, ok := t.Params[k]
	if !ok {
		return 0, errors.New("no such parameter: " + k)
	}
	f, ok := i.(float64)
	if !ok {
		return 0, fmt.Errorf("%s is not a number: %T", k, i)
	}
	return f, nil
}

// GetAsInt converts a number parameter associated with k in t.Params into int.
func (t *PluginParams) GetAsInt(k string) (int, error) {
	f, err := t.GetAsFloat(k)
	if err != nil {
		return 0, err
	}
	return int(f), nil
}

// GetAsBool returns a boolean parameter associated with k in t.Params.
func (t *PluginParams) GetAsBool(k string) (bool, error) {
	i, ok := t.Params[k]
	if !ok {
		return false, errors.New("no such parameter: " + k)
	}
	b, ok := i.(bool)
	if !ok {
		return false, fmt.Errorf("%s is not a number: %T", k, i)
	}
	return b, nil
}
