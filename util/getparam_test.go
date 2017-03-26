package util

import (
	"math"
	"reflect"
	"testing"
)

func TestGetString(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{
		"k1": 1,
		"k2": nil,
		"k3": "abc",
	}

	if _, err := GetString("k", m); !IsNotFound(err) {
		t.Error(`_, err := GetString(m, "k"); !IsNotFound(err)`)
	}

	if _, err := GetString("k1", m); !IsBadType(err) {
		t.Error(`_, err := GetString(m, "k1"); !IsBadType(err)`)
	}

	if s, err := GetString("k2", m); err != nil {
		t.Error(err)
	} else if s != "" {
		t.Error(`s != ""`)
	}

	if s, err := GetString("k3", m); err != nil {
		t.Error(err)
	} else if s != "abc" {
		t.Error(`s != "abc"`)
	}
}

func TestGetBool(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{
		"k1": 1,
		"k2": nil,
		"k3": true,
		"k4": false,
	}

	if _, err := GetBool("k", m); !IsNotFound(err) {
		t.Error(`_, err := GetBool(m, "k"); !IsNotFound(err)`)
	}

	if _, err := GetBool("k1", m); !IsBadType(err) {
		t.Error(`_, err := GetBool(m, "k1"); !IsBadType(err)`)
	}

	if s, err := GetBool("k2", m); err != nil {
		t.Error(err)
	} else if s != false {
		t.Error(`s != false`)
	}

	if s, err := GetBool("k3", m); err != nil {
		t.Error(err)
	} else if s != true {
		t.Error(`s != true`)
	}

	if s, err := GetBool("k4", m); err != nil {
		t.Error(err)
	} else if s != false {
		t.Error(`s != false`)
	}
}

func TestGetInt(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{
		"k1": "abc",
		"k2": nil,
		"k3": 10,
		"k4": 11.0,
		"k5": 11.1,
		"k6": math.NaN(),
		"k7": math.Inf(0),
		"k8": -2,
		"k9": float64(-5.0),
	}

	if _, err := GetInt("k", m); !IsNotFound(err) {
		t.Error(`_, err := GetInt("k", m); !IsNotFound(err)`)
	}

	if _, err := GetInt("k1", m); !IsBadType(err) {
		t.Error(`_, err := GetInt("k1", m); !IsBadType(err)`)
	}

	if _, err := GetInt("k2", m); !IsBadType(err) {
		t.Error(`_, err := GetInt("k2", m); !IsBadType(err)`)
	}

	if v, err := GetInt("k3", m); err != nil {
		t.Error(err)
	} else if v != 10 {
		t.Error(`v != 10`)
	}

	if v, err := GetInt("k4", m); err != nil {
		t.Error(err)
	} else if v != 11 {
		t.Error(`v != 11`)
	}

	if _, err := GetInt("k5", m); !IsBadType(err) {
		t.Error(`_, err := GetInt("k5", m); !IsBadType(err)`)
	}

	if _, err := GetInt("k6", m); !IsBadType(err) {
		t.Error(`_, err := GetInt("k6", m); !IsBadType(err)`)
	}

	if _, err := GetInt("k7", m); !IsBadType(err) {
		t.Error(`_, err := GetInt("k7", m); !IsBadType(err)`)
	}

	if v, err := GetInt("k8", m); err != nil {
		t.Error(err)
	} else if v != -2 {
		t.Error(`v != -2`)
	}

	if v, err := GetInt("k9", m); err != nil {
		t.Error(err)
	} else if v != -5 {
		t.Error(`v != -5`)
	}
}

func TestGetFloat64(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{
		"k1": "abc",
		"k2": nil,
		"k3": 10,
		"k4": 11.0,
		"k5": 11.1,
		"k6": math.NaN(),
		"k7": math.Inf(0),
		"k8": -2,
		"k9": float64(-5.0),
	}

	if _, err := GetFloat64("k", m); !IsNotFound(err) {
		t.Error(`_, err := GetFloat64("k", m); !IsNotFound(err)`)
	}

	if _, err := GetFloat64("k1", m); !IsBadType(err) {
		t.Error(`_, err := GetFloat64("k1", m); !IsBadType(err)`)
	}

	if _, err := GetFloat64("k2", m); !IsBadType(err) {
		t.Error(`_, err := GetFloat64("k2", m); !IsBadType(err)`)
	}

	if v, err := GetFloat64("k3", m); err != nil {
		t.Error(err)
	} else if v != float64(10) {
		t.Error(`v != float64(10)`)
	}

	if v, err := GetFloat64("k4", m); err != nil {
		t.Error(err)
	} else if v != 11.0 {
		t.Error(`v != 11.0`)
	}

	if v, err := GetFloat64("k5", m); err != nil {
		t.Error(err)
	} else if v != 11.1 {
		t.Error(`v != 11.1`)
	}

	if v, err := GetFloat64("k6", m); err != nil {
		t.Error(err)
	} else if !math.IsNaN(v) {
		t.Error(`!math.IsNaN(v)`)
	}

	if v, err := GetFloat64("k7", m); err != nil {
		t.Error(err)
	} else if !math.IsInf(v, 0) {
		t.Error(`!math.IsInf(v)`)
	}

	if v, err := GetFloat64("k8", m); err != nil {
		t.Error(err)
	} else if v != float64(-2) {
		t.Error(`v != float64(-2)`)
	}

	if v, err := GetFloat64("k9", m); err != nil {
		t.Error(err)
	} else if v != -5.0 {
		t.Error(`v != -5.0`)
	}
}

func TestGetStringSlice(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{
		"k1": 1,
		"k2": nil,
		"k3": []interface{}{1, "abc"},
		"k4": []interface{}{"abc", "def"},
		"k5": []string{"abc", "def"},
	}

	if _, err := GetStringSlice("k", m); !IsNotFound(err) {
		t.Error(`_, err := GetStringSlice("k", m); !IsNotFound(err)`)
	}

	if _, err := GetStringSlice("k1", m); !IsBadType(err) {
		t.Error(`_, err := GetStringSlice("k1", m); !IsBadType(err)`)
	}

	if v, err := GetStringSlice("k2", m); err != nil {
		t.Error(err)
	} else if len(v) != 0 {
		t.Error(`len(v) != 0`)
	}

	if _, err := GetStringSlice("k3", m); !IsBadType(err) {
		t.Error(`_, err := GetStringSlice("k3", m); !IsBadType(err)`)
	}

	if v, err := GetStringSlice("k4", m); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, []string{"abc", "def"}) {
		t.Error(`!reflect.DeepEqual(v, []string{"abc", "def"})`)
	}

	if v, err := GetStringSlice("k5", m); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, []string{"abc", "def"}) {
		t.Error(`!reflect.DeepEqual(v, []string{"abc", "def"})`)
	}
}

func TestGetIntSlice(t *testing.T) {
	t.Parallel()

	m := map[string]interface{}{
		"k1": "abc",
		"k2": nil,
		"k3": []interface{}{1, "abc"},
		"k4": []interface{}{1, 2},
		"k5": []interface{}{1, float64(2.0)},
		"k6": []interface{}{1, 2.1},
		"k7": []int{1, 2},
		"k8": []float64{1.0, 2.0},
	}

	if _, err := GetIntSlice("k", m); !IsNotFound(err) {
		t.Error(`_, err := GetIntSlice("k", m); !IsNotFound(err)`)
	}

	if _, err := GetIntSlice("k1", m); !IsBadType(err) {
		t.Error(`_, err := GetIntSlice("k1", m); !IsBadType(err)`)
	}

	if v, err := GetIntSlice("k2", m); err != nil {
		t.Error(err)
	} else if len(v) != 0 {
		t.Error(`len(v) != 0`)
	}

	if _, err := GetIntSlice("k3", m); !IsBadType(err) {
		t.Error(`_, err := GetIntSlice("k3", m); !IsBadType(err)`)
	}

	if v, err := GetIntSlice("k4", m); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, []int{1, 2}) {
		t.Error(`!reflect.DeepEqual(v, []int{1, 2})`)
	}

	if v, err := GetIntSlice("k5", m); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, []int{1, 2}) {
		t.Error(`!reflect.DeepEqual(v, []int{1, 2})`)
	}

	if _, err := GetIntSlice("k6", m); !IsBadType(err) {
		t.Error(`_, err := GetIntSlice("k6", m); !IsBadType(err)`)
	}

	if v, err := GetIntSlice("k7", m); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(v, []int{1, 2}) {
		t.Error(`!reflect.DeepEqual(v, []int{1, 2})`)
	}

	if _, err := GetIntSlice("k8", m); !IsBadType(err) {
		t.Error(`_, err := GetIntSlice("k8", m); !IsBadType(err)`)
	}
}
