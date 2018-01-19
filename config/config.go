/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package config

import (
	"gopkg.in/yaml.v2"
	"time"
)

//=================================================================================================
// build vars
//=================================================================================================

var (
	Codename     = "Project Title"
	Version      = "0.0.1"
	Commit       = "dev"
	BuildTime    = time.Now().Format(time.UnixDate)
	BuilderName  = "builder"
	BuildMachine = "localhost"
)

//=================================================================================================
// definition of config structure
//=================================================================================================

type Config struct {
	Node struct {
		Name string `yaml:"name"`
		Dc   string `yaml:"dc"`
		Rack string `yaml:"rack"`
	} `yaml:"node"`

	Runtime struct {
		MaxProcs int `yaml:"maxprocs"`
	} `yaml:"runtime"`

	Dirs struct {
		Data string `yaml:"data"`
		Tmp  string `yaml:"tmp"`
	} `yaml:"dirs"`

	Log struct {
		Destination string `yaml:"destination"`
		Level       string `yaml:"level"`
	} `yaml:"log"`

	Metrics struct {
		Destination      string `yaml:"destination"`
		UpdateIntervalMs int    `yaml:"interval_ms"`
	} `yaml:"metrics"`

	DataSources struct {
		Db struct {
			ConnectionString       string `yaml:"connect"`
			PoolMaxIdleConnections int    `yaml:"pool_max_idle_connections"`
			PoolMaxConnections     int    `yaml:"pool_max_connections"`
		} `yaml:"db"`
	} `yaml:"data_sources"`

	Services struct {
		Acounter struct {
			Enabled bool      `yaml:"enabled"`
			Tag     string    `yaml:"tag"`
			Socket  SrvSocket `yaml:"socket"`
		}
	} `yaml:"services"`
}

type SrvSocket struct {
	Proto string `yaml:"proto"`
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Path  string `yaml:"path"`
}

//=================================================================================================
// helpers
//=================================================================================================

func Parse(configPayload string) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal([]byte(configPayload), config)
	if err != nil {
		return config, err
	}
	return config, nil
}
