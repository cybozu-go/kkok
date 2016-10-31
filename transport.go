package kkok

import "errors"

// Transport is the interface that transport plugins must implement.
type Transport interface {

	// Params returns PluginParams that can be used to construct
	// this transport.
	Params() PluginParams

	// String should return a descriptive one-line string of the transport.
	String() string

	// Deliver alerts via the transport.
	//
	// The transport may merge alerts into one, or may deliver
	// alerts one by one.
	Deliver(alerts []*Alert) error
}

// TransportConstructor is a function signature for transport construction.
type TransportConstructor func(params map[string]interface{}) (Transport, error)

var transportTypes = make(map[string]TransportConstructor)

// RegisterTransport registers a construction function of a Transport type.
func RegisterTransport(typ string, ctor TransportConstructor) {
	transportTypes[typ] = ctor
}

// NewTransport constructs a Transport.
func NewTransport(typ string, params map[string]interface{}) (Transport, error) {
	ctor, ok := transportTypes[typ]
	if !ok {
		return nil, errors.New("no such transport type: " + typ)
	}
	return ctor(params)
}
