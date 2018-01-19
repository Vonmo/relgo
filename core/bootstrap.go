/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */
package core

import (
	"fmt"
	"github.com/Vonmo/relgo/config"
	"github.com/Vonmo/relgo/log"
	"github.com/Vonmo/relgo/metrics"
	"github.com/Vonmo/relgo/models"
	"io/ioutil"
	"os"
	"os/signal"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

func (core *Core) bootstrap(configFile string) {
	core.Interrupt = make(chan os.Signal, 1)
	core.Lock = &sync.RWMutex{}
	core.Services = make(map[string]*Service)

	core.initLog(configFile)
	core.initMetrics()
	core.initRuntime()
	core.initDirs()
	core.initDb()
	core.handleTermSignal()

	// core is ready.
	log.Debug("start ", System.Config.Node.Name, "...")
	core.Ready = true
}

func (core *Core) initRuntime() {
	runtime.GOMAXPROCS(core.Config.Runtime.MaxProcs)
}

func (core *Core) initDb() {
	models.InitDB(
		core.Config.DataSources.Db.ConnectionString,
		core.Config.DataSources.Db.PoolMaxConnections,
		core.Config.DataSources.Db.PoolMaxIdleConnections,
		core.Metrics,
	)
}

func (core *Core) handleTermSignal() {
	go func() {
		<-core.Interrupt
		log.Print("received term signal")
		for _, s := range core.Services {
			if s.ShutdownFun != nil {
				s.ShutdownFun("TERM")
			}
		}
		os.Exit(0)
	}()
	signal.Notify(core.Interrupt, os.Interrupt, syscall.SIGTERM)
}

func (core *Core) initLog(configFile string) {
	payloadBytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Panicf("config issue: %s", err.Error())
	}
	core.Config, err = config.Parse(string(payloadBytes))
	if err != nil {
		log.Fatal(err)
	}
	log.Init(core.Config.Log.Level, true, core.Config.Log.Destination)
	log.Debugf("config file: " + configFile)
}

func (core *Core) initMetrics() {
	core.Metrics = metrics.Init(
		core.Config.Metrics.Destination,
		core.Config.Metrics.UpdateIntervalMs,
	)
	core.Metrics.Set("start_time", time.Now().UTC().String())
	core.Metrics.Set("nodename", core.Config.Node.Name)
	core.Metrics.Set("dc", core.Config.Node.Dc)
	core.Metrics.Set("dc_rack", core.Config.Node.Rack)
	core.Metrics.Set("runtime.max_proc", int64(core.Config.Runtime.MaxProcs))
}

func (core *Core) initDirs() {
	s := reflect.ValueOf(&core.Config.Dirs).Elem()
	for i := 0; i < s.NumField(); i++ {
		log.Debugf("ensure %s dir: %s", s.Type().Field(i).Name, s.Field(i).Interface())
		ensureDir(s.Field(i).Interface().(string), os.ModePerm)
	}
}

func ensureDir(path string, perm os.FileMode) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, perm)
	}
}

func projectTitle() (title string) {
	return fmt.Sprintf(
		"\n%s.\n%s-%s %s (%s@%s)\n%s\n",
		config.Codename, config.Version, config.Commit, config.BuildTime,
		config.BuilderName, config.BuildMachine, copyright(),
	)
}

func PrintBuildInfo() {
	fmt.Printf(
		"Programm Name:\t%s\nVersion:\t%s-%s\nBuild date:\t%s\nMaintainer:\t%s\nBuild hostname:\t%s\nCopyright:\t%s\n",
		config.Codename, config.Version, config.Commit, config.BuildTime,
		config.BuilderName, config.BuildMachine, copyright(),
	)
}

func copyright() (copyright string) {
	return "(c) " + strconv.Itoa(time.Now().Year()) + " Vonmo."
}
