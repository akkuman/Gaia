package rdp

import (
	"Gaia/plugin"
	"sync"
)

var pluginName = "rdp"

// BurstPlugin ftp burst plugin
type BurstPlugin struct {
}

// Flag turn on?
func (BurstPlugin) Flag() bool {
	return plugin.SwitchIsOn(pluginName)
}

// Start start a burst
func (BurstPlugin) Start(wgMain *sync.WaitGroup) {
	plugin.TaskCellStart(pluginName, wgMain, threadFunc)
}

func threadFunc(info [3]string) (success bool, err error) {
	return
}

func init() {
	ftpPlugin := BurstPlugin{}
	plugin.Register(pluginName, ftpPlugin)
}
