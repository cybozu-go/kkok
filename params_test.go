package kkok

import (
	"encoding/json"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/cybozu-go/kkok/util"
)

func testTOML(t *testing.T) {
	t.Parallel()

	d := `
type = "fuga"
p1 = 3
p2 = 1.234
`
	var p PluginParams
	_, err := toml.Decode(d, &p)
	if err != nil {
		t.Error(err)
	}

	if p.Type != "fuga" {
		t.Error(`p.Type != "fuga"`)
	}
	if p.Params["p1"].(int64) != 3 {
		t.Error(`p.Params["p1"] != 3`)
	}
	p2 := p.Params["p2"].(float64)
	if p2 < 1.2 || p2 > 1.3 {
		t.Error(`p2 < 1.2 || p2 > 1.3`)
	}
}

func testJSON(t *testing.T) {
	t.Parallel()

	p := PluginParams{
		Type: "fuga",
		Params: map[string]interface{}{
			"p1": 3,
			"p2": 1.234,
			"p3": true,
		},
	}

	j, err := json.Marshal(&p)
	if err != nil {
		t.Fatal(err)
	}

	var pp PluginParams
	err = json.Unmarshal(j, &pp)
	if err != nil {
		t.Fatal(err)
	}

	if pp.Type != "fuga" {
		t.Error(`pp.Type != "fuga"`)
	}
	if p1, err := util.GetInt("p1", pp.Params); err != nil {
		t.Error(err)
	} else if p1 != 3 {
		t.Error(`p1 != 3`)
	}
	if p2, err := util.GetFloat64("p2", pp.Params); err != nil {
		t.Error(err)
	} else if p2 < 1.2 || p2 > 1.3 {
		t.Error(`p2 < 1.2 || p2 > 1.3`)
	}
	if p3, err := util.GetBool("p3", pp.Params); err != nil {
		t.Error(err)
	} else if !p3 {
		t.Error(`!p3`)
	}
}

func TestPluginParams(t *testing.T) {
	t.Run("TOML", testTOML)
	t.Run("JSON", testJSON)
}
