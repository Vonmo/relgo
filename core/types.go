/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */
package core

import (
	"github.com/Vonmo/relgo/config"
	"github.com/Vonmo/relgo/metrics"
	"github.com/gobuffalo/packr"
	"os"
	"sync"
)

type CoreOptions struct {
	ConfigFile  string
	CommitedVsn string
	Boxes       *CoreBoxes
}

type CoreBoxes struct {
	Static packr.Box
}

type Service struct {
	Name        string
	Ready       bool
	ShutdownFun func(reason string)
	Dict        map[string]interface{}
}

type Core struct {
	Services map[string]*Service
	Config   *config.Config
	Metrics  *metrics.Metrics
	Boxes    *CoreBoxes
	Ready    bool

	Lock      *sync.RWMutex
	Interrupt chan os.Signal
}
