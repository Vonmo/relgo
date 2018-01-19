/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */
package main

import (
	"flag"
	"github.com/Vonmo/relgo/core"
	"github.com/Vonmo/relgo/log"
	"github.com/Vonmo/relgo/services"
	"github.com/gobuffalo/packr"
	"os"
)

var (
	configFile    = flag.String("config", "config.yml", "Node Config")
	showBuildInfo = flag.Bool("V", false, "Show build information")
)

func main() {
	defer func() {
		str := recover()
		if str != nil {
			log.Panic(str)
			os.Exit(1)
		} else {
			log.Print(services.Engine + " done.")
			os.Exit(0)
		}
	}()

	flag.Parse()

	if *showBuildInfo {
		core.PrintBuildInfo()
		os.Exit(0)
	}

	log.Fatal(core.Run(&core.CoreOptions{
		ConfigFile: *configFile,
		Boxes: &core.CoreBoxes{
			Static: packr.NewBox("./static"),
		},
	}))
}
