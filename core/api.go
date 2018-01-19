/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */
package core

import (
	"fmt"
	"github.com/Vonmo/relgo/config"
	"time"
)

var System *Core

func Run(options *CoreOptions) (result string) {
	if config.Commit != "dev" {
		fmt.Println(projectTitle())
	}

	System = &Core{Boxes: options.Boxes}
	System.bootstrap(options.ConfigFile)

	for {
		System.Metrics.Set("services.count", int64(System.countServices()))
		time.Sleep(1 * time.Second)
	}

	return "core was stopped"
}

func Register(srv *Service) {
	System.registerService(srv)
}
