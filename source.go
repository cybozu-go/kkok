package kkok

import (
	"context"
	"errors"
)

// Source is the interface that generators must implement.
type Source interface {

	// Run the generator until ctx.Done() is canceled.
	// Use post to post the generated alerts.
	Run(ctx context.Context, post func(*Alert)) error
}

// SourceConstructor is a function signature for source construction.
type SourceConstructor func(params map[string]interface{}) (Source, error)

var sourceTypes = make(map[string]SourceConstructor)

// RegisterSource registers a construction function of a Source type.
func RegisterSource(typ string, ctor SourceConstructor) {
	sourceTypes[typ] = ctor
}

// NewSource constructs a Source.
func NewSource(typ string, params map[string]interface{}) (Source, error) {
	ctor, ok := sourceTypes[typ]
	if !ok {
		return nil, errors.New("no such source type: " + typ)
	}
	return ctor(params)
}
