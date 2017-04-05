package freq

import (
	"errors"

	"github.com/cybozu-go/kkok"
)

type clType int

const (
	clNone clType = iota
	clFrom
	clTitle
	clHost
)

func str2cl(s string) (clType, error) {
	switch s {
	case "From", "from":
		return clFrom, nil
	case "Title", "title":
		return clTitle, nil
	case "Host", "host":
		return clHost, nil
	}
	return clNone, errors.New("no such class: " + s)
}

// Value returns the key string according to the classify type.
func (cl clType) Value(a *kkok.Alert) string {
	switch cl {
	case clFrom:
		return a.From
	case clTitle:
		return a.Title
	case clHost:
		return a.Host
	}
	return ""
}

// String returns a string for classify type.
func (cl clType) String() string {
	switch cl {
	case clFrom:
		return "From"
	case clTitle:
		return "Title"
	case clHost:
		return "Host"
	}
	return ""
}
