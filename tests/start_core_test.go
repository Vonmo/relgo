package tests

import (
	"github.com/elzor/relgo/core"
	"github.com/gobuffalo/packr"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	go core.Run(&core.CoreOptions{
		ConfigFile: "../config/config.yml",
		Boxes: &core.CoreBoxes{
			Static: packr.NewBox("./static"),
		},
	})
	os.Exit(m.Run())
}
