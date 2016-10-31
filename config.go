package kkok

import (
	"time"

	"github.com/cybozu-go/cmd"
)

const (
	defaultInitialInterval = 30
	defaultMaxInterval     = 30
	defaultAddr            = ":19898"
)

// Config is a struct to load TOML configuration file for kkok.
type Config struct {
	// InitialInterval specifies the initial interval seconds
	// to pool posted alerts before procession.
	// Default is 30 (seconds).
	InitialInterval int `toml:"initial_interval"`

	// MaxInterval specifies the maximum interval seconds to
	// pool posted alerts.  The interval begins with InitialInterval
	// then doubles if one or more alerts are posted during the
	// interval until it reaches MaxInterval.
	//
	// The interval will return to InitialInterval if no alerts are
	// posted during the current interval.
	//
	// Default is 30 (seconds).
	MaxInterval int `toml:"max_interval"`

	// Addr is the listen address for HTTP API.
	//
	// Default is ":19898"
	Addr string `toml:"listen"`

	// APIToken is used for API authentication if not empty.
	//
	// Default is empty.
	APIToken string `toml:"api_token"`

	// Log from cybozu-go/cmd.
	Log cmd.LogConfig `toml:"log"`

	// Sources is a list of parameters to construct alert generators.
	Sources []PluginParams `toml:"source"`

	// Routes is a map between route ID and a list of transports.
	Routes map[string][]PluginParams `toml:"route"`

	// Filters is a list of parameters to construct filters.
	Filters []PluginParams `toml:"filter"`
}

// InitialDuration returns the initial dispatch interval.
func (c *Config) InitialDuration() time.Duration {
	return time.Second * time.Duration(c.InitialInterval)
}

// MaxDuration returns the maximum dispatch interval.
func (c *Config) MaxDuration() time.Duration {
	return time.Second * time.Duration(c.MaxInterval)
}

// NewConfig returns *Config with default settings.
func NewConfig() *Config {
	return &Config{
		InitialInterval: defaultInitialInterval,
		MaxInterval:     defaultMaxInterval,
		Addr:            defaultAddr,
	}
}
