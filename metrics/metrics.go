/*
 * Created: 2017.06
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */

package metrics

import (
	"github.com/Vonmo/relgo/log"
	"io/ioutil"
	"reflect"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Metrics struct {
	data  map[string]interface{}
	mutex *sync.RWMutex
}

func Init(outputPath string, updateIntervalMs int) *Metrics {
	log.Debug("metrics path: " + outputPath)

	metrics := &Metrics{
		data:  make(map[string]interface{}),
		mutex: &sync.RWMutex{},
	}

	go func() {
		defer func() {
			str := recover()
			if str != nil {
				log.Error(str)
				Init(outputPath, updateIntervalMs)
			}
		}()
		for {

			metrics.Set("current_time", time.Now().UTC().String())

			keys := []string{}
			metrics.mutex.RLock()
			for k, _ := range metrics.data {
				keys = append(keys, k)
			}
			metrics.mutex.RUnlock()
			sort.Strings(keys)

			text := ""
			for _, k := range keys {
				text += k + ": "
				metrics.mutex.RLock()
				if reflect.TypeOf(metrics.data[k]).Name() == "string" {
					text += metrics.data[k].(string)
				} else {
					text += strconv.FormatInt(metrics.data[k].(int64), 10)
				}
				metrics.mutex.RUnlock()
				text += "\n"
			}
			err := ioutil.WriteFile(outputPath, []byte(text), 0644)
			if err != nil {
				log.Error(err)
			}
			time.Sleep(time.Duration(updateIntervalMs) * time.Millisecond)
		}
	}()

	return metrics
}

func (metrics *Metrics) Increase(name string) {
	go func() {
		defer metrics.mutex.Unlock()
		for {
			metrics.mutex.Lock()
			if val, ok := metrics.data[name]; ok {
				metrics.data[name] = val.(int64) + 1
			} else {
				metrics.data[name] = int64(1)
			}
			break
		}
	}()
}

func (metrics *Metrics) Decrease(name string) {
	go func() {
		defer metrics.mutex.Unlock()
		for {
			metrics.mutex.Lock()
			if val, ok := metrics.data[name]; ok {
				metrics.data[name] = val.(int64) - 1
			} else {
				metrics.data[name] = int64(-1)
			}
			break
		}
	}()
}

func (metrics *Metrics) Set(name string, val interface{}) {
	go func() {
		defer metrics.mutex.Unlock()
		for {
			metrics.mutex.Lock()
			delete(metrics.data, name)
			metrics.data[name] = val
			break
		}
	}()
}

func (metrics *Metrics) Get(name string) interface{} {
	defer metrics.mutex.RUnlock()
	var val interface{}
	for {
		metrics.mutex.RLock()
		if v, ok := metrics.data[name]; ok {
			val = v
		}
		break
	}
	return val
}
