/*
 * Created: 2018.01
 * Author: Maxim Molchanov <m.molchanov@vonmo.com>
 */
package core

import "time"

func NewService(name string) Service {
	return Service{
		Name:  name,
		Ready: false,
		Dict:  map[string]interface{}{},
	}
}

func (core *Core) countServices() (cnt int) {
	for {
		core.Lock.RLock()
		cnt = len(core.Services)
		core.Lock.RUnlock()
		break
	}
	return cnt
}

func (core *Core) registerService(srv *Service) {
	for {
		core.Lock.Lock()
		core.Services[srv.Name] = srv
		core.Lock.Unlock()
		break
	}
}

func WaitCore() {
	for {
		if System != nil && System.Ready {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func WaitService(srv string) {
	for {
		if System != nil && System.Ready {
			if _, ok := System.Services[srv]; ok {
				if System.Services[srv].Ready {
					break
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}
