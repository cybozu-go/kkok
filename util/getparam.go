package util

import (
	"errors"
	"math"
)

var (
	errNotFound = errors.New("not found")
	errBadType  = errors.New("bad type")
)

// IsNotFound returns true if err indicates the specified key was not found.
func IsNotFound(err error) bool {
	return err == errNotFound
}

// IsBadType returns true if err indicates the type of value is bad.
func IsBadType(err error) bool {
	return err == errBadType
}

// GetString looks up m[k] and return the value if it is a string.
// As a special case, if m[k] is nil, this returns "".
//
// If err is not nil, it can be tested with IsNotFound or IsBadType.
func GetString(k string, m map[string]interface{}) (v string, err error) {
	i, ok := m[k]
	if !ok {
		err = errNotFound
		return
	}

	switch i := i.(type) {
	case nil:
	case string:
		v = i
	default:
		err = errBadType
	}
	return
}

// GetBool looks up m[k] and return the value if it is a boolean.
// As a special case, if m[k] is nil, this returns false.
//
// If err is not nil, it can be tested with IsNotFound or IsBadType.
func GetBool(k string, m map[string]interface{}) (v bool, err error) {
	i, ok := m[k]
	if !ok {
		err = errNotFound
		return
	}

	switch i := i.(type) {
	case nil:
	case bool:
		v = i
	default:
		err = errBadType
	}
	return
}

// GetInt looks up m[k] and return the value if it is an int.
// As a special case, if the value is float64 and
// math.Trunc(v) == v, it returns int(math.Trunc(v)).
//
// If err is not nil, it can be tested with IsNotFound or IsBadType.
func GetInt(k string, m map[string]interface{}) (v int, err error) {
	i, ok := m[k]
	if !ok {
		err = errNotFound
		return
	}

	switch i := i.(type) {
	case int:
		v = i
	case float64:
		switch {
		case math.IsNaN(i), math.IsInf(i, 0):
			err = errBadType
		case math.Trunc(i) == i:
			v = int(i)
		default:
			err = errBadType
		}
	default:
		err = errBadType
	}
	return
}

// GetFloat64 looks up m[k] and return the value if it is an int or float64.
//
// If err is not nil, it can be tested with IsNotFound or IsBadType.
func GetFloat64(k string, m map[string]interface{}) (v float64, err error) {
	i, ok := m[k]
	if !ok {
		err = errNotFound
		return
	}

	switch i := i.(type) {
	case int:
		v = float64(i)
	case float64:
		v = i
	default:
		err = errBadType
	}
	return
}

// GetStringSlice looks up m[k] and return the value
// if it is a slice of strings.
// As a special case, if m[k] is nil, this returns nil.
//
// If err is not nil, it can be tested with IsNotFound or IsBadType.
func GetStringSlice(k string, m map[string]interface{}) (v []string, err error) {
	i, ok := m[k]
	if !ok {
		err = errNotFound
		return
	}

	switch i := i.(type) {
	case nil:
	case []interface{}:
		if len(i) == 0 {
			return
		}
		sl := make([]string, len(i))
		for idx, ii := range i {
			s, ok := ii.(string)
			if !ok {
				err = errBadType
				return
			}
			sl[idx] = s
		}
		v = sl
	case []string:
		v = i
	default:
		err = errBadType
	}
	return
}

// GetIntSlice looks up m[k] and return the value
// if it is a slice of int or float64.
// As a special case, if m[k] is nil, this returns nil.
//
// If err is not nil, it can be tested with IsNotFound or IsBadType.
func GetIntSlice(k string, m map[string]interface{}) (v []int, err error) {
	i, ok := m[k]
	if !ok {
		err = errNotFound
		return
	}

	switch i := i.(type) {
	case nil:
	case []interface{}:
		if len(i) == 0 {
			return
		}
		il := make([]int, len(i))
		for idx, ii := range i {
			switch ii := ii.(type) {
			case int:
				il[idx] = ii
			case float64:
				switch {
				case math.IsNaN(ii), math.IsInf(ii, 0):
					err = errBadType
					return
				case math.Trunc(ii) == ii:
					il[idx] = int(ii)
				default:
					err = errBadType
					return
				}
			default:
				err = errBadType
				return
			}
		}
		v = il
	case []int:
		v = i
	default:
		err = errBadType
	}
	return
}
