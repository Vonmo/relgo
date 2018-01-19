package services

import (
	"github.com/elzor/relgo/core"
	"time"
)

const Engine = "VonmoSRVS"

func waitCore() {
	for {
		if core.System != nil && core.System.Ready {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}
